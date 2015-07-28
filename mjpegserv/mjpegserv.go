package mjpegserv

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

// MJPEGServCreator is a creator of MJPEG server.
type MJPEGServCreator struct{}

// CreateSink creates a MJPEG server sink, user can access AVI file or images
// through HTTP access.
//
// Usage of WITH parameters:
//  server_name: [required] a server name
//  port:        a port number, default port is 10090
//
// Example:
//  when a creation query is
//    `CREATE SINK hoge_result TYPE mjpeg_server WITH server_name='foo', port=8080`
//  then the sink addressed http://localhost:8080/video/foo/
func (m *MJPEGServCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	servName, err := params.Get("server_name")
	if err != nil {
		return nil, fmt.Errorf("mjpeg server requires the server name")
	}
	servNameStr, err := data.AsString(servName)
	if err != nil {
		return nil, fmt.Errorf("mjpeg server name needs to define as string type")
	}

	port, err := params.Get("port")
	if err != nil {
		defaultPort := 10090
		ctx.Log().Infof("mjpeg server is starting with %d port", defaultPort)
		port = data.Int(defaultPort)
	}
	portNum, err := data.AsInt(port)
	if err != nil {
		return nil, err
	}

	ms := &mjpegServ{}
	ms.serverName = servNameStr
	ms.port = int(portNum)
	ms.inChan = make(chan input)
	go ms.start()
	return ms, nil
}

func (m *MJPEGServCreator) TypeName() string {
	return "mjpeg_server"
}

type inputData struct {
	name      string
	imageData []byte
}

type input struct {
	key       string
	inputData inputData
}

type mjpegServ struct {
	serverName string
	port       int
	pub        *publisher
	inChan     chan input
}

func (m *mjpegServ) Write(ctx *core.Context, t *core.Tuple) error {
	name, err := t.Data.Get("name")
	if err != nil {
		return err
	}
	nameStr, err := data.AsString(name)
	if err != nil {
		return err
	}

	img, err := t.Data.Get("img")
	if err != nil {
		return err
	}
	imgByte, err := data.AsBlob(img)
	if err != nil {
		return err
	}
	imgp := bridge.DeserializeMatVec3b(imgByte)
	defer imgp.Delete()

	inData := inputData{
		name:      nameStr,
		imageData: imgp.ToJpegData(50),
	}
	in := input{
		key:       m.serverName,
		inputData: inData,
	}

	m.inChan <- in
	return nil
}

func (m *mjpegServ) Close(ctx *core.Context) error {
	// closing web server is better
	return nil
}

type subscriber struct {
	id  int
	key string
	ch  chan []byte
}

func (s *subscriber) channel() chan []byte {
	return s.ch
}

type publisher struct {
	subscribers     map[int]*subscriber
	nextSubsriberId int
	m               sync.Mutex
	currentFrames   map[string][]byte
}

func newPublisher() *publisher {
	return &publisher{
		subscribers:   map[int]*subscriber{},
		currentFrames: map[string][]byte{},
	}
}

func (p *publisher) getNameList() []string {
	names := []string{}
	for name, _ := range p.currentFrames {
		names = append(names, name)
	}
	return names
}

func (p *publisher) updateFrame(key string, f []byte) {
	p.m.Lock()
	p.currentFrames[key] = f
	p.m.Unlock()
	go p.publish(key)
}

func (p *publisher) publish(key string) {
	p.m.Lock()
	defer p.m.Unlock()
	for _, s := range p.subscribers {
		if s.key != key {
			continue
		}
		select {
		case s.ch <- p.currentFrames[key]:
		default:
		}
	}
}

func (p *publisher) subscribe(key string) *subscriber {
	p.m.Lock()
	defer p.m.Unlock()
	s := &subscriber{
		id:  p.nextSubsriberId,
		key: key,
		ch:  make(chan []byte),
	}
	p.nextSubsriberId += 1
	go func() {
		s.ch <- p.currentFrames[key]
	}()
	p.subscribers[s.id] = s
	return s
}

func (p *publisher) close(s *subscriber) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.subscribers, s.id)
}

func (m *mjpegServ) start() {
	pub := newPublisher()

	fin := make(chan bool)
	go func() {
		for {
			select {
			case in := <-m.inChan:
				go pub.updateFrame(in.inputData.name, in.inputData.imageData)
			default:
			}
		}
	}()

	go func() {
		router := web.New(mpegServContext{
			pub: pub,
		})
		router.Get(`/video/:name`, (*mpegServContext).videoHandler) // TODO regex validation
		router.Get(`/snapshot/:key`, (*mpegServContext).snapshotHandler)
		router.Get(`/list`, (*mpegServContext).listHandler)
		if err := http.ListenAndServe(fmt.Sprint(":", m.port), router); err != nil {
			log.Println("cannot start the server: ", err)
		}
		fin <- true
	}()

	<-fin
}

type mpegServContext struct {
	pub *publisher
}

func (m *mpegServContext) videoHandler(rw web.ResponseWriter, req *web.Request) {
	log.Println(req.URL.Path)
	name, ok := req.PathParams["name"]
	if !ok {
		log.Println("Not found: ", req.URL)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	conn, bufrw, err := rw.Hijack()
	if err != nil {
		log.Println("Failed to hijack a connection: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	log.Println("Video streaming connection started: ", req.RemoteAddr)

	s := m.pub.subscribe(name)
	defer m.pub.close(s)
	log.Println("Started to subscribe ", name)

	bufrw.WriteString("HTTP/1.1 200 OK\r\n")
	msg := "Content-Type: multipart/x-mixed-replace; boundary=\"myboundary\"\r\n\r\n--myboundary"
	if _, err := bufrw.WriteString(msg); err != nil {
		log.Println("Failed to write header")
		return
	}
	bufrw.Flush()

	for f := range s.channel() {
		head := fmt.Sprintf("\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(f))
		bufrw.WriteString(head)
		bufrw.Write(f)
		_, err = bufrw.WriteString("\r\n--myboundary")
		if err != nil {
			log.Println("Write failed")
			return
		}
		bufrw.Flush()
	}
	bufrw.WriteString("--")
	bufrw.Flush()
}

func (m *mpegServContext) snapshotHandler(rw web.ResponseWriter, req *web.Request) {
	log.Println(req.URL.Path)
	key, ok := req.PathParams["key"]
	if !ok {
		log.Println("Not found: ", req.URL)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	s := m.pub.subscribe(key)
	defer m.pub.close(s)

	rw.Header().Set("Content-Type", "image/jpeg")
	rw.Write(<-s.channel())
}

func (m *mpegServContext) listHandler(rw web.ResponseWriter, req *web.Request) {
	nameList := m.pub.getNameList()

	rw.Header().Set("Content-Type", "text/html")
	for _, name := range nameList {
		rw.Write([]byte(fmt.Sprintf("<a href='/video/%s'><img src='/video/%s' title='%s'></a>\n", name, name, name)))
	}
}

package writer

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
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
//  port:    A port number, default port is 10090.
//  quality: The quality of converting JPEG file, if empty then set 50.
//
// Example:
//  when a creation query is
//    `CREATE SINK mjpeg_moniter TYPE scouter_mjpeg_server WITH port=8080`
//  then the sink addressed http://localhost:8080/video/foo/ (Address of "foo"
//  will be decided by `INSERT` query.)
//
// A created AVI file can not be opened until this sink is dropped or SensorBee
// process is down.
func (m *MJPEGServCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	port, err := params.Get(utils.PortPath)
	if err != nil {
		defaultPort := 10090
		ctx.Log().Infof("mjpeg server is starting with %d port", defaultPort)
		port = data.Int(defaultPort)
	}
	portNum, err := data.AsInt(port)
	if err != nil {
		return nil, err
	}

	quality, err := params.Get(utils.QualityPath)
	if err != nil {
		quality = data.Int(50)
	}
	q, err := data.AsInt(quality)
	if err != nil {
		return nil, err
	}

	ms := &mjpegServ{}
	ms.finish = false
	ms.port = int(portNum)
	ms.quality = int(q)
	ms.pub = newPublisher()
	ms.inChan = make(chan inputData)
	go ms.start()
	return ms, nil
}

// TypeName returns type name.
func (m *MJPEGServCreator) TypeName() string {
	return "scouter_mjpeg_server"
}

type inputData struct {
	name      string
	imageData []byte
}

type mjpegServ struct {
	finish  bool
	port    int
	quality int
	pub     *publisher
	inChan  chan inputData
}

// Write input images to a server which have started when sink creation.
// Input tuple is required to have follow `data.Map`
//
//  data.Map{
//    "name": [access category name] (will be casted to string type)
//    "img" : [image binary data] (`data.Blob`)
//  }
//
// Example of insertion query:
//  ```
//  INSERT INTO mjpeg_monitor SELECT ISTREAM
//    detection_ressult_frame AS img,
//    `camera1_detection` AS name
//    FROM detected_frame @RANGE 1 TUPLES];
//  ```
// then URI will be
//  * http://localhot:8080/video/camera1_detection
//    Users can watch images, which updated automatically.
//  * http://localhot:8080/snapshot/camera1_detection
//    Users can see a snapshot image.
func (m *mjpegServ) Write(ctx *core.Context, t *core.Tuple) error {
	name, err := t.Data.Get(utils.NamePath)
	if err != nil {
		return err
	}
	nameStr, err := data.AsString(name)
	if err != nil {
		return err
	}

	img, err := t.Data.Get(utils.IMGPath)
	if err != nil {
		return err
	}
	imgByte, err := data.ToBlob(img)
	if err != nil {
		return err
	}
	imgp := bridge.DeserializeMatVec3b(imgByte)
	defer imgp.Delete()

	data := inputData{
		name:      nameStr,
		imageData: imgp.ToJpegData(m.quality),
	}

	m.inChan <- data
	return nil
}

func (m *mjpegServ) Close(ctx *core.Context) error {
	// closing web server is better
	m.finish = true
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
	nextSubsriberID int
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
	for name := range p.currentFrames {
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
		id:  p.nextSubsriberID,
		key: key,
		ch:  make(chan []byte),
	}
	p.nextSubsriberID++
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
	fin := make(chan bool)
	go func() {
		for !m.finish {
			select {
			case in := <-m.inChan:
				go m.pub.updateFrame(in.name, in.imageData)
			default:
			}
		}
	}()

	go func() {
		ctx := mpegServContext{
			pub: m.pub,
		}
		router := web.New(ctx)
		// TODO regex validation
		router.Get(`/video/:name`, ctx.videoHandler)
		router.Get(`/snapshot/:key`, ctx.snapshotHandler)
		router.Get(`/list`, ctx.listHandler)
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
	msg := "Content-Type: multipart/x-mixed-replace; "
	msg += "boundary=\"myboundary\"\r\n\r\n--myboundary"
	if _, err := bufrw.WriteString(msg); err != nil {
		log.Println("Failed to write header")
		return
	}
	bufrw.Flush()

	for f := range s.channel() {
		head := fmt.Sprintf(
			"\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(f))
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
	formatTag := "<a href='/video/%s'><img src='/video/%s' title='%s'></a>\n"
	for _, name := range nameList {
		rw.Write([]byte(fmt.Sprintf(formatTag, name, name, name)))
	}
}

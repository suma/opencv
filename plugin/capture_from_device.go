package plugin

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"strconv"
	"sync"
	"time"
)

// CaptureFromDevice is a frame generator using OpenCV video capture.
// Usage of WITH parameters:
//  DeviceID: [required] the ID of associated device
//  Width: frame width, if set empty or "0" then will be ignore
//  Height: frame height, if set empty or "0" then will be ignore
//  FPS: frame per second, if set empty or "0" then will be ignore
//  CameraID: the unique ID of this source if set empty then the ID will be 0
type CaptureFromDevice struct {
	vcap   bridge.VideoCapture
	finish bool

	DeviceID int64
	Width    int64
	Height   int64
	FPS      int64
	CameraID int64
}

// GenerateStream streams video capture datum. OpenCV parameters
// (e.g width, height...) are set in struct members.
func (c *CaptureFromDevice) GenerateStream(ctx *core.Context, w core.Writer) error {
	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.OpenDevice(int(c.DeviceID)); !ok {
		return fmt.Errorf("error opening device: %v", c.DeviceID)
	}

	// OpenCV video capture configuration
	if c.Width > 0 {
		c.vcap.Set(conf.CvCapPropFrameWidth, int(c.Width))
	}
	if c.Height > 0 {
		c.vcap.Set(conf.CvCapPropFrameHeight, int(c.Height))
	}
	if c.FPS > 0 {
		c.vcap.Set(conf.CvCapPropFps, int(c.FPS))
	}

	// read camera frames
	mu := sync.RWMutex{}
	rootBuf := bridge.NewMatVec3b()
	defer rootBuf.Delete()
	var rootBufErr error
	errChan := make(chan error)
	go func(rootBuf bridge.MatVec3b) {
		for {
			c.grab(rootBuf, &mu, errChan)
			select {
			case err := <-errChan:
				if err != nil {
					rootBufErr = err
					break
				}
			}
		}
	}(rootBuf)

	// streaming, capture from rootBuf
	buf := bridge.NewMatVec3b()
	defer buf.Delete()
	c.finish = false
	for !c.finish {
		if rootBufErr != nil {
			return rootBufErr
		}
		mu.RLock()
		rootBuf.CopyTo(buf)
		mu.RUnlock()
		if buf.Empty() {
			continue
		}

		var m = tuple.Map{
			"capture":  tuple.Blob(buf.Serialize()),
			"cameraID": tuple.Int(c.CameraID),
		}
		now := time.Now()
		t := tuple.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now,
			Trace:         make([]tuple.TraceEvent, 0),
		}
		w.Write(ctx, &t)
	}
	return nil
}

func (c *CaptureFromDevice) grab(buf bridge.MatVec3b, mu *sync.RWMutex, errChan chan error) {
	if !c.vcap.IsOpened() {
		errChan <- fmt.Errorf("video stream or file closed, device no: %d)", c.DeviceID)
		return
	}
	tmpBuf := bridge.NewMatVec3b()
	defer tmpBuf.Delete()
	if ok := c.vcap.Read(tmpBuf); !ok {
		errChan <- fmt.Errorf("cannot read a new file (device no: %d)", c.DeviceID)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	tmpBuf.CopyTo(buf)
}

func (c *CaptureFromDevice) Stop(ctx *core.Context) error {
	c.finish = true
	c.vcap.Delete()
	return nil
}

func (c *CaptureFromDevice) GetSourceCreator() (bql.SourceCreator, error) {
	creator := func(with map[string]string) (core.Source, error) {
		did, ok := with["device_id"]
		if !ok {
			return nil, fmt.Errorf("capture source need device ID")
		}
		deviceID, err := strconv.ParseInt(did, 10, 64)
		if err != nil {
			return nil, err
		}

		w, ok := with["width"]
		if !ok {
			w = "0" // will be ignored
		}
		width, err := strconv.ParseInt(w, 10, 64)
		if err != nil {
			return nil, err
		}

		h, ok := with["height"]
		if !ok {
			h = "0" // will be ignored
		}
		height, err := strconv.ParseInt(h, 10, 64)
		if err != nil {
			return nil, err
		}

		f, ok := with["fps"]
		if !ok {
			f = "0" // will be ignored
		}
		fps, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, err
		}

		cid, ok := with["camera_id"]
		if !ok {
			cid = "0"
		}
		cameraID, err := strconv.ParseInt(cid, 10, 64)
		if err != nil {
			return nil, err
		}

		c.DeviceID = deviceID
		c.Width = width
		c.Height = height
		c.FPS = fps
		c.CameraID = cameraID
		return c, nil
	}
	return creator, nil
}

func (c *CaptureFromDevice) TypeName() string {
	return "capture_from_device"
}

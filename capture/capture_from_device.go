package capture

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
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
		c.vcap.Set(bridge.CvCapPropFrameWidth, int(c.Width))
	}
	if c.Height > 0 {
		c.vcap.Set(bridge.CvCapPropFrameHeight, int(c.Height))
	}
	if c.FPS > 0 {
		c.vcap.Set(bridge.CvCapPropFps, int(c.FPS))
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

func (c *CaptureFromDevice) CreateSource(ctx *core.Context, with tuple.Map) (core.Source, error) {
	did, err := with.Get("device_id")
	if err != nil {
		return nil, err
	}
	deviceID, err := tuple.AsInt(did)
	if err != nil {
		return nil, err
	}

	w, err := with.Get("width")
	if err != nil {
		w = tuple.Int(0) // will be ignored
	}
	width, err := tuple.AsInt(w)
	if err != nil {
		return nil, err
	}

	h, err := with.Get("height")
	if err != nil {
		h = tuple.Int(0) // will be ignored
	}
	height, err := tuple.AsInt(h)
	if err != nil {
		return nil, err
	}

	f, err := with.Get("fps")
	if err != nil {
		f = tuple.Int(0) // will be ignored
	}
	fps, err := tuple.AsInt(f)
	if err != nil {
		return nil, err
	}

	cid, err := with.Get("camera_id")
	if err != nil {
		cid = tuple.Int(0)
	}
	cameraID, err := tuple.AsInt(cid)
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

func (c *CaptureFromDevice) TypeName() string {
	return "capture_from_device"
}

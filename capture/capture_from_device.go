package capture

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"time"
)

// CaptureFromDeviceCreator is a creator of a capture from device.
type CaptureFromDeviceCreator struct{}

func (c *CaptureFromDeviceCreator) TypeName() string {
	return "capture_from_device"
}

// CreateSource creates a frame generator using OpenCV video capture.
//
// Usage of WITH parameters:
//  device_id: [required] the ID of associated device
//  width:     frame width, if set empty or "0" then will be ignore
//  height:    frame height, if set empty or "0" then will be ignore
//  fps:       frame per second, if set empty or "0" then will be ignore
//  camera_id: the unique ID of this source if set empty then the ID will be 0
func (c *CaptureFromDeviceCreator) CreateSource(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Source, error) {
	did, err := params.Get("device_id")
	if err != nil {
		return nil, err
	}
	deviceID, err := data.AsInt(did)
	if err != nil {
		return nil, err
	}

	w, err := params.Get("width")
	if err != nil {
		w = data.Int(0) // will be ignored
	}
	width, err := data.AsInt(w)
	if err != nil {
		return nil, err
	}

	h, err := params.Get("height")
	if err != nil {
		h = data.Int(0) // will be ignored
	}
	height, err := data.AsInt(h)
	if err != nil {
		return nil, err
	}

	f, err := params.Get("fps")
	if err != nil {
		f = data.Int(0) // will be ignored
	}
	fps, err := data.AsInt(f)
	if err != nil {
		return nil, err
	}

	cid, err := params.Get("camera_id")
	if err != nil {
		cid = data.Int(0)
	}
	cameraID, err := data.AsInt(cid)
	if err != nil {
		return nil, err
	}

	cs := &captureFromDevice{}
	cs.deviceID = deviceID
	cs.width = width
	cs.height = height
	cs.fps = fps
	cs.cameraID = cameraID
	return cs, nil
}

type captureFromDevice struct {
	vcap   bridge.VideoCapture
	finish bool

	deviceID int64
	width    int64
	height   int64
	fps      int64
	cameraID int64
}

// GenerateStream streams video capture datum. OpenCV parameters
// (e.g width, height...) are set when the source is initialized.
func (c *captureFromDevice) GenerateStream(ctx *core.Context, w core.Writer) error {
	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.OpenDevice(int(c.deviceID)); !ok {
		return fmt.Errorf("error opening device: %v", c.deviceID)
	}

	// OpenCV video capture configuration
	if c.width > 0 {
		c.vcap.Set(bridge.CvCapPropFrameWidth, int(c.width))
	}
	if c.height > 0 {
		c.vcap.Set(bridge.CvCapPropFrameHeight, int(c.height))
	}
	if c.fps > 0 {
		c.vcap.Set(bridge.CvCapPropFps, int(c.fps))
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

		now := time.Now()
		var m = data.Map{
			"capture":   data.Blob(buf.Serialize()),
			"cameraID":  data.Int(c.cameraID),
			"timestamp": data.Timestamp(now),
		}
		t := core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now,
			Trace:         []core.TraceEvent{},
		}
		w.Write(ctx, &t)
	}
	return nil
}

func (c *captureFromDevice) grab(buf bridge.MatVec3b, mu *sync.RWMutex, errChan chan error) {
	if !c.vcap.IsOpened() {
		errChan <- fmt.Errorf("video stream or file closed, device no: %d)", c.deviceID)
		return
	}
	tmpBuf := bridge.NewMatVec3b()
	defer tmpBuf.Delete()
	if ok := c.vcap.Read(tmpBuf); !ok {
		errChan <- fmt.Errorf("cannot read a new file (device no: %d)", c.deviceID)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	tmpBuf.CopyTo(buf)
}

func (c *captureFromDevice) Stop(ctx *core.Context) error {
	c.finish = true
	c.vcap.Delete()
	return nil
}

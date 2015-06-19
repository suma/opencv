package plugin

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
	"time"
)

// CaptureFromDevice is a frame generator using OpenCV video capture.
type CaptureFromDevice struct {
	vcap   bridge.VideoCapture
	finish bool

	DeviceID int
	Width    int
	Height   int
	FPS      int
	CameraId int
}

// GenerateStream streams video capture datum. OpenCV parameters
// (e.g width, height...) are set in struct members.
func (c *CaptureFromDevice) GenerateStream(ctx *core.Context, w core.Writer) error {
	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.OpenDevice(c.DeviceID); !ok {
		return fmt.Errorf("error opening device: %v", c.DeviceID)
	}

	// OpenCV video capture configuration
	if c.Width > 0 {
		c.vcap.Set(conf.CvCapPropFrameWidth, c.Width)
	}
	if c.Height > 0 {
		c.vcap.Set(conf.CvCapPropFrameHeight, c.Height)
	}
	if c.FPS > 0 {
		c.vcap.Set(conf.CvCapPropFps, c.FPS)
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
			"cameraID": tuple.Int(c.CameraId),
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

func (c *CaptureFromDevice) Schema() *core.Schema {
	return nil
}

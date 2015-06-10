package snippets

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	devicePrefix = "device://"
)

// Capture reads camera frame or video data using OpenCV library
type Capture struct {
	config conf.CaptureConfig
	vcap   bridge.VideoCapture
	fp     bridge.FrameProcessor
	finish bool
}

// SetUp prepares Video Capture. Video information and reading frame rate
// is set up by external configuration file.
func (c *Capture) SetUp(configFilePath string) error {
	config, err := conf.GetCaptureSnippetConfig(configFilePath)
	if err != nil {
		return err
	}
	c.config = config
	c.vcap = bridge.NewVideoCapture()

	if strings.HasPrefix(config.URI, devicePrefix) {
		deviceNoStr := config.URI[len(devicePrefix):len(config.URI)]
		deviceNo, err := strconv.Atoi(deviceNoStr)
		if err != nil {
			return fmt.Errorf("error opening device: %v", deviceNoStr)
		}
		if ok := c.vcap.OpenDevice(deviceNo); !ok {
			return fmt.Errorf("error opening device: %v", deviceNoStr)
		}
		if config.Width > 0 {
			c.vcap.Set(conf.CvCapPropFrameWidth, config.Width)
		}
		if config.Height > 0 {
			c.vcap.Set(conf.CvCapPropFrameHeight, config.Height)
		}
		if config.TickInterval > 0 {
			c.vcap.Set(conf.CvCapPropFps, 1000.0/config.TickInterval)
		}
	} else {
		if ok := c.vcap.Open(config.URI); !ok {
			return fmt.Errorf("error opening video stream or file: %v", config.URI)
		}
	}

	c.fp = bridge.NewFrameProcessor(config.FrameProcessorConfig)
	c.finish = false

	return nil
}

func grab(vcap bridge.VideoCapture, buf bridge.MatVec3b, mu *sync.RWMutex, errChan chan error) {
	if !vcap.IsOpened() {
		errChan <- fmt.Errorf("video stream or file closed")
		return
	}
	tmpBuf := bridge.NewMatVec3b()
	defer tmpBuf.Delete()
	if ok := vcap.Read(tmpBuf); !ok {
		errChan <- fmt.Errorf("cannot read a new frame")
		return
	}

	mu.Lock()
	defer mu.Unlock()
	tmpBuf.CopyTo(buf)
}

func (c *Capture) stopFlag() {
	c.finish = true
}

// GenerateStream generates tuples include captured frame. If an error
// occur, this streaming stop.
func (c *Capture) GenerateStream(ctx *core.Context, w core.Writer) error {
	mu := sync.RWMutex{}
	defer c.stopFlag()

	config := c.config
	rootBuf := bridge.NewMatVec3b()
	defer rootBuf.Delete()
	var rootBufErr error
	if !config.CaptureFromFile {
		errChan := make(chan error)
		go func(rootBuf bridge.MatVec3b) {
			for {
				grab(c.vcap, rootBuf, &mu, errChan)
				select {
				case err := <-errChan:
					if err != nil {
						rootBufErr = err
						break
					}
				}
			}
		}(rootBuf)
	}

	buf := bridge.NewMatVec3b()
	defer buf.Delete()
	cnt := 0
	for !c.finish {
		cnt++
		if config.CaptureFromFile {
			if ok := c.vcap.Read(buf); !ok {
				ctx.Logger.Log(core.Debug, "total read frames count is %d", cnt)
				return fmt.Errorf("cannot read a new frame")
			}
			if config.FrameSkip > 0 {
				c.vcap.Grab(config.FrameSkip)
			}
		} else {
			if rootBufErr != nil {
				return rootBufErr
			}

			mu.RLock()
			rootBuf.CopyTo(buf)
			mu.RUnlock()
			if buf.Empty() {
				continue
			}
		}

		now := time.Now()
		inow := now.UnixNano() / int64(time.Millisecond) // [ms]
		f := c.fp.Apply(buf, inow, config.CameraID)

		var m = tuple.Map{
			"frame":     tuple.Blob(f.Serialize()),
			"camera_id": tuple.Int(config.CameraID),
		}
		t := tuple.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now, // TODO video capture create time
			Trace:         make([]tuple.TraceEvent, 0),
		}
		w.Write(ctx, &t)
	}
	return nil
}

// Stop generating stream.
func (c *Capture) Stop(ctx *core.Context) error {
	c.finish = true
	time.Sleep(500 * time.Millisecond)
	c.fp.Delete()
	c.vcap.Delete()
	return nil
}

// Schema returns registered schema.
func (c *Capture) Schema() *core.Schema {
	return nil
}

// IsStopped returns whether this capture has stoppoed or not.
// If the capture is stopped, return true.
func (c *Capture) IsStopped() bool {
	return c.finish
}

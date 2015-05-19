package snippets

import (
	"C"
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type CaptureConfig struct {
	CameraID        int
	URI             string
	CaptureFromFile bool
	FrameSkip       int
}

type Capture struct {
	config CaptureConfig
	vcap   bridge.VideoCapture
	fp     bridge.FrameProcessor
}

func (c *Capture) SetUp(config CaptureConfig) error {
	c.config = config
	var vcap bridge.VideoCapture
	ok := bridge.VideoCapture_Open(config.URI, vcap)
	if !ok {
		return fmt.Errorf("error opening video stream or file : %v", config.URI)
	}
	c.vcap = vcap

	var fp bridge.FrameProcessor
	bridge.FrameProcessor_SetUp(fp, nil)
	c.fp = fp

	return nil
}

func grab(vcap bridge.VideoCapture, buf bridge.MatVec3b, errChan chan error) {
	if !bridge.VideoCapture_IsOpened(vcap) {
		errChan <- fmt.Errorf("video stream or file closed")
		return
	}
	var tmpBuf bridge.MatVec3b
	ok := bridge.VideoCapture_Read(vcap, tmpBuf)
	if !ok {
		errChan <- fmt.Errorf("cannot read a new frame")
		return
	}
	bridge.MatVec3b_Clone(tmpBuf, buf)
}

func (c *Capture) GenerateStream(ctx *core.Context, w core.Writer) error {
	config := c.config
	var rootBuf, buf bridge.MatVec3b
	var rootBufErr error
	if !config.CaptureFromFile {
		errChan := make(chan error)
		go func(rootBuf bridge.MatVec3b) {
			for {
				grab(c.vcap, rootBuf, errChan)
				select {
				case err := <-errChan:
					rootBufErr = err
					break
				}
			}
		}(rootBuf)
	}

	for { // TODO add stop command and using goroutine
		if config.CaptureFromFile {
			ok := bridge.VideoCapture_Read(c.vcap, buf)
			if !ok {
				return fmt.Errorf("cannot read a new frame")
			}
			if config.FrameSkip > 0 {
				for i := 0; i < config.FrameSkip; i++ {
					// TODO considering biding cost
					bridge.VideoCapture_Grab(c.vcap)
				}
			}
		} else {
			if rootBufErr != nil {
				return rootBufErr
			}
			bridge.MatVec3b_Clone(rootBuf, buf)
			if bridge.MatVec3b_Empty(buf) {
				continue
			}
		}

		// TODO create frame processor configuration, very difficult...
		// TODO confirm time stamp using, create in C++ is better?
		now := time.Now()
		inow, _ := tuple.ToInt(tuple.Timestamp(now))
		_, f := bridge.FrameProcessor_Apply(c.fp, buf, inow, config.CameraID)

		var m = tuple.Map{
			"frame": tuple.Blob(f),
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

func (c *Capture) Stop(ctx *core.Context) error {
	return nil
}

func (c *Capture) Schema() *core.Schema {
	return nil
}

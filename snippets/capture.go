package snippets

import (
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
}

func (c *Capture) SetUp(config CaptureConfig) error {
	c.config = config
	vcap := bridge.VideoCapture_Open(config.URI)
	if vcap == nil {
		return fmt.Errorf("error opening video stream or file : %v", config.URI)
	}
	c.vcap = vcap

	return nil
}

func (c *Capture) GenerateStream(ctx *core.Context, w core.Writer) error {
	var buf bridge.MatVec3b
	config := c.config
	for { // TODO add stop command and using goroutine
		if config.CaptureFromFile {
			bridge.VideoCapture_Read(c.vcap, buf)
			if buf == nil {
				return fmt.Errorf("cannot read a new frame")
			}
			if config.FrameSkip > 0 {
				for i := 0; i < config.FrameSkip; i++ {
					// TODO considering biding cost
					bridge.VideoCapture_Grab(c.vcap)
				}
			}
		} else {
			buf = bridge.MatVec3b_Clone(buf)
			if bridge.MatVec3b_Empty(buf) {
				continue
			}
		}

		// TODO create frame processor configuration, very difficult...
		// TODO confirm time stamp using, create in C++ is better?
		now := time.Now()
		inow, _ := tuple.ToInt(tuple.Timestamp(now))
		fp := bridge.FrameProcessor_SetUp(nil)
		f := bridge.FrameProcessor_Apply(fp, buf, inow, config.CameraID)

		var m = tuple.Map{
			"frame": tuple.Blob(f),
		}
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

func (c *Capture) Schema() *core.Schema {
	return nil
}

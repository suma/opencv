package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type CaptureConfig struct {
	CameraId         int
	Uri              string
	CaptureFromMFile bool
	FrameSkip        int
}

type Capture struct {
	config CaptureConfig
	vcap   bridge.VideoCapture
}

func (c *Capture) SetUp(config CaptureConfig) error {
	c.config = config
	vcap := bridge.VideoCapture_Open(config.Uri)
	if vcap == nil {
		return fmt.Errorf("error opening video stream or file : %v", config.Uri)
	}
	c.vcap = vcap

	return nil
}

func (c *Capture) GenerateStream(ctx *core.Context, w core.Writer) error {
	var buf bridge.MatVec3b
	config := c.config
	for { // TODO add stop command
		if config.CaptureFromMFile {
			bridge.VideoCapture_Read(c.vcap, buf)
			if buf == nil {
				return fmt.Errorf("cannot read a new frame")
			}
			if config.FrameSkip > 0 {
				for i := 0; i < config.FrameSkip; i++ {
					bridge.VideoCapture_Grab(c.vcap)
				}
			}
		} else {
			buf = bridge.MatVec3b_Clone(buf)
			if bridge.MatVec3b_Empty(buf) {
				continue
			}
		}
	}
	// TOBE get frames
	now := time.Now()
	t := tuple.Tuple{
		Timestamp:     now,
		ProcTimestamp: now,
		Trace:         make([]tuple.TraceEvent, 0),
	}
	w.Write(ctx, &t)
	return nil
}

func (c *Capture) Schema() *core.Schema {
	return nil
}

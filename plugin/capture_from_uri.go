package plugin

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

// CaptureFromURI is a frame generator using OpenCV video capture.
// URI can be set URL and file path.
type CaptureFromURI struct {
	vcap   bridge.VideoCapture
	finish bool

	URI       string
	FrameSkip int
	CameraId  int
}

// GenerateStream streams video capture datum. OpenCV video capture read
// frames from URI, but user can control frame streaming frequency use
// FrameSkip.
//
// !ATTENTION!
// When a capture source is a file-style (e.g. AVI file) and complete to read
// all frames, an error will be occurred because video capture cannot read
// a new frame. User can count total frame count to confirm complete of read
// file. The number of count is logged.
func (c *CaptureFromURI) GenerateStream(ctx *core.Context, w core.Writer) error {
	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.Open(c.URI); !ok {
		return fmt.Errorf("error opening video stream or file: %v", c.URI)
	}

	buf := bridge.NewMatVec3b()
	defer buf.Delete()
	cnt := 0
	c.finish = false
	for !c.finish {
		cnt++
		if ok := c.vcap.Read(buf); !ok {
			ctx.Logger.Log(core.Debug, "total read frames count is %d", cnt-1)
			return fmt.Errorf("cannot read a new frame: %v", c.URI)
		}
		if c.FrameSkip > 0 {
			c.vcap.Grab(c.FrameSkip)
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

func (c *CaptureFromURI) Stop(ctx *core.Context) error {
	c.finish = true
	c.vcap.Delete()
	return nil
}

func (c *CaptureFromURI) Schema() *core.Schema {
	return nil
}

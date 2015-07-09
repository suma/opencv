package capture

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

// CaptureFromURI is a frame generator using OpenCV video capture.
// URI can be set HTTP address or file path.
// Usage of WITH parameters:
//  URI: [required] a capture data's URI (e.g. /data/test.avi)
//  FrameSkip: the number of frame skip, if set empty or "0" then read all frames
//  CameraID: the unique ID of this source if set empty then the ID will be 0
type CaptureFromURI struct {
	vcap   bridge.VideoCapture
	finish bool

	URI       string
	FrameSkip int64
	CameraID  int64
}

// GenerateStream streams video capture datum. OpenCV video capture read
// frames from URI, user can control frame streaming frequency use
// FrameSkip.
//
// !ATTENTION!
// When a capture source is a file-style (e.g. AVI file) and complete to read
// all frames, an error will be occurred because video capture cannot read
// a new frame. User can count total frame to confirm complete of read file.
// The number of count is logged.
func (c *CaptureFromURI) GenerateStream(ctx *core.Context, w core.Writer) error {
	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.Open(c.URI); !ok {
		return fmt.Errorf("error opening video stream or file: %v", c.URI)
	}

	buf := bridge.NewMatVec3b()
	defer buf.Delete()
	cnt := 0
	c.finish = false
	ctx.Logger.Log(core.Debug, "start reading video stream of file: %v", c.URI)
	for !c.finish {
		cnt++
		if ok := c.vcap.Read(buf); !ok {
			ctx.Logger.Log(core.Debug, "total read frames count is %d", cnt-1)
			return fmt.Errorf("cannot read a new frame: %v", c.URI)
		}
		if c.FrameSkip > 0 {
			c.vcap.Grab(int(c.FrameSkip))
		}

		var m = data.Map{
			"capture":   data.Blob(buf.Serialize()),
			"camera_id": data.Int(c.CameraID),
		}
		now := time.Now()
		t := core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now,
			Trace:         make([]core.TraceEvent, 0),
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

func (c *CaptureFromURI) CreateSource(ctx *core.Context, with data.Map) (core.Source, error) {
	uri, err := with.Get("uri")
	if err != nil {
		return nil, fmt.Errorf("capture source needs URI")
	}
	uriStr, err := data.AsString(uri)
	if err != nil {
		return nil, err
	}

	fs, err := with.Get("frame_skip")
	if err != nil {
		fs = data.Int(0) // will be ignored
	}
	frameSkip, err := data.AsInt(fs)
	if err != nil {
		return nil, err
	}

	cid, err := with.Get("camera_id")
	if err != nil {
		cid = data.Int(0)
	}
	cameraID, err := data.AsInt(cid)
	if err != nil {
		return nil, err
	}

	c.URI = uriStr
	c.FrameSkip = frameSkip
	c.CameraID = cameraID
	return c, nil
}

func (c *CaptureFromURI) TypeName() string {
	return "capture_from_uri"
}

package capture

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync/atomic"
	"time"
)

// CaptureFromURICreator is a creator of a capture from URI.
type CaptureFromURICreator struct{}

func (c *CaptureFromURICreator) TypeName() string {
	return "capture_from_uri"
}

// CreateSource creates a frame generator using OpenCV video capture.
// URI can be set HTTP address or file path.
//
// Usage of WITH parameters:
//  uri:        [required] a capture data's URI (e.g. /data/test.avi)
//  frame_skip: the number of frame skip, if set empty or "0" then read all frames
//  camera_id:  the unique ID of this source if set empty then the ID will be 0
func (c *CaptureFromURICreator) CreateSource(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Source, error) {

	cs, err := createCaptureFromURI(ctx, ioParams, params)
	if err != nil {
		return nil, err
	}
	//return core.NewRewindableSource(cs), nil
	return cs, nil // TODO GenerateStream cannot be called concurrently
}

func createCaptureFromURI(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (
	core.Source, error) {

	uri, err := params.Get("uri")
	if err != nil {
		return nil, fmt.Errorf("capture source needs URI")
	}
	uriStr, err := data.AsString(uri)
	if err != nil {
		return nil, err
	}

	fs, err := params.Get("frame_skip")
	if err != nil {
		fs = data.Int(0) // will be ignored
	}
	frameSkip, err := data.AsInt(fs)
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

	cs := &captureFromURI{}
	cs.finish = false
	cs.paused = 0
	cs.uri = uriStr
	cs.frameSkip = frameSkip
	cs.cameraID = cameraID
	return cs, nil
}

type captureFromURI struct {
	vcap   bridge.VideoCapture
	finish bool
	// paused is used as atomic bool
	// paused set 0 then means false, set other then means true
	paused int32

	uri       string
	frameSkip int64
	cameraID  int64
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
func (c *captureFromURI) GenerateStream(ctx *core.Context, w core.Writer) error {

	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.Open(c.uri); !ok {
		return fmt.Errorf("error opening video stream or file: %v", c.uri)
	}

	buf := bridge.NewMatVec3b()
	defer buf.Delete()
	cnt := 0
	c.finish = false
	ctx.Log().Infof("start reading video stream of file: %v", c.uri)
	for !c.finish {
		if atomic.LoadInt32(&(c.paused)) != 0 {
			continue
		}
		cnt++
		if ok := c.vcap.Read(buf); !ok {
			ctx.Log().Infof("total read frames count is %d", cnt-1)
			atomic.StoreInt32(&(c.paused), int32(1))
			continue
		}
		if c.frameSkip > 0 {
			c.vcap.Grab(int(c.frameSkip))
		}

		var m = data.Map{
			"capture":   data.Blob(buf.Serialize()),
			"camera_id": data.Int(c.cameraID),
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

func (c *captureFromURI) Stop(ctx *core.Context) error {
	c.finish = true
	c.vcap.Delete()
	return nil
}

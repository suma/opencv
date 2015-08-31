package capture

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync/atomic"
	"time"
)

// FromURICreator is a creator of a capture from URI.
type FromURICreator struct{}

// TypeName returns type name.
func (c *FromURICreator) TypeName() string {
	return "scouter_capture_from_uri"
}

// CreateSource creates a frame generator using OpenCV video capture.
// URI can be set HTTP address or file path.
//
// Usage of WITH parameters:
//  uri:              [required] A capture data's URI (e.g. /data/test.avi).
//  frame_skip:       The number of frame skip, if set empty or "0" then read
//                    all frames. FPS is depended on the URI's file (or device).
//  camera_id:        The unique ID of this source if set empty then the ID will
//                    be 0.
//  next_frame_error: When this source cannot read a new frame, occur error or
//                    not decided by the flag. If the flag set `true` then
//                    return error. Default value is true.
//  rewind:           If set `true` then user can use `REWIND SOURCE` query.
func (c *FromURICreator) CreateSource(ctx *core.Context,
	ioParams *bql.IOParams, params data.Map) (core.Source, error) {

	cs, err := createCaptureFromURI(ctx, ioParams, params)
	if err != nil {
		return nil, err
	}

	rewindFlag := false
	if rf, err := params.Get(utils.RewindPath); err == nil {
		if rewindFlag, err = data.AsBool(rf); err != nil {
			return nil, err
		}
	}
	if rewindFlag {
		return core.NewRewindableSource(cs), nil
	}
	return cs, nil
}

func createCaptureFromURI(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (
	core.Source, error) {

	uri, err := params.Get(utils.URIPath)
	if err != nil {
		return nil, fmt.Errorf("capture source needs URI")
	}
	uriStr, err := data.AsString(uri)
	if err != nil {
		return nil, err
	}

	fs, err := params.Get(utils.FrameSkipPath)
	if err != nil {
		fs = data.Int(0) // will be ignored
	}
	frameSkip, err := data.AsInt(fs)
	if err != nil {
		return nil, err
	}

	cid, err := params.Get(utils.CameraIDPath)
	if err != nil {
		cid = data.Int(0)
	}
	cameraID, err := data.AsInt(cid)
	if err != nil {
		return nil, err
	}

	endErrFlag, err := params.Get(utils.NextFrameErrorPath)
	if err != nil {
		endErrFlag = data.True
	}
	endErr, err := data.AsBool(endErrFlag)
	if err != nil {
		return nil, err
	}

	cs := &captureFromURI{}
	atomic.StoreInt32(&(cs.stop), int32(1))
	cs.uri = uriStr
	cs.frameSkip = frameSkip
	cs.cameraID = cameraID
	cs.endErrFlag = endErr
	return cs, nil
}

type captureFromURI struct {
	vcap bridge.VideoCapture
	// stop is used as atomic bool
	// stop set 0 then means false, set other then means true
	stop int32

	uri        string
	frameSkip  int64
	cameraID   int64
	endErrFlag bool
}

// GenerateStream streams video capture datum. OpenCV video capture read
// frames from URI, user can control frame streaming frequency using
// FrameSkip.This source is rewindable.
//
// Output:
//  capture:   The frame image binary data ('data.Blob'), serialized from
//             OpenCV's matrix data format (`cv::Mat_<cv::Vec3b>`).
//  camera_id: The camera ID.
//  timestamp: The timestamp of capturing. (reed below details)
//
// When a capture source is a file-style (e.g. AVI file), a "timestamp" value is
// NOT correspond with the file created time. The "timestamp" value is the time
// of this source capturing a new frame.
// And when complete to read the file's all frames, video capture cannot read a
// new frame. If the key "next_frame_error" set `false` then a no new frame
// error will not be occurred, User can also count the number of total frame to
// confirm complete of read file. The number of count is logged.
func (c *captureFromURI) GenerateStream(ctx *core.Context, w core.Writer) error {
	if atomic.LoadInt32(&(c.stop)) == 0 {
		atomic.StoreInt32(&(c.stop), int32(1))
		ctx.Log().Infof("interrupt reading video stream or file and reset: %v",
			c.uri)
		c.vcap.Release()
		c.vcap.Delete()
	}

	c.vcap = bridge.NewVideoCapture()
	if ok := c.vcap.Open(c.uri); !ok {
		return fmt.Errorf("error opening video stream or file: %v", c.uri)
	}

	buf := bridge.NewMatVec3b()
	defer buf.Delete()

	cnt := 0
	ctx.Log().Infof("start reading video stream of file: %v", c.uri)
	atomic.StoreInt32(&(c.stop), int32(0))
	defer atomic.StoreInt32(&(c.stop), int32(1))
	for atomic.LoadInt32(&(c.stop)) == 0 {
		cnt++
		if ok := c.vcap.Read(buf); !ok {
			ctx.Log().Infof("total read frames count is %d", cnt-1)
			if c.endErrFlag {
				return fmt.Errorf("cannot reed a new frame")
			}
			break
		}
		if c.frameSkip > 0 {
			c.vcap.Grab(int(c.frameSkip))
		}

		now := time.Now()
		var m = data.Map{
			"capture":   data.Blob(buf.Serialize()),
			"camera_id": data.Int(c.cameraID),
			"timestamp": data.Timestamp(now),
		}
		t := core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now,
			Trace:         []core.TraceEvent{},
		}
		err := w.Write(ctx, &t)
		if err == core.ErrSourceRewound || err == core.ErrSourceStopped {
			return err
		}
	}
	return nil
}

func (c *captureFromURI) Stop(ctx *core.Context) error {
	atomic.StoreInt32(&(c.stop), int32(1))
	c.vcap.Delete()
	return nil
}

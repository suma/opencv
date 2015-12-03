package opencv

import (
	"fmt"
	"pfi/sensorbee/opencv/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"time"
)

// FromDeviceCreator is a creator of a capture from device.
type FromDeviceCreator struct{}

var (
	deviceIDPath = data.MustCompilePath("device_id")
	widthPath    = data.MustCompilePath("width")
	heightPath   = data.MustCompilePath("height")
	fpsPath      = data.MustCompilePath("fps")
)

// CreateSource creates a frame generator using OpenCV video capture
// (`VideoCapture::open`).
//
// Usage of WITH parameters:
//  device_id: [required] The ID of associated device.
//  width:     Frame width, if set empty or "0" then will be ignore.
//  height:    Frame height, if set empty or "0" then will be ignore.
//  fps:       Frame per second, if set empty or "0" then will be ignore.
//  camera_id: The unique ID of this source if set empty then the ID will be 0.
func (c *FromDeviceCreator) CreateSource(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Source, error) {
	did, err := params.Get(deviceIDPath)
	if err != nil {
		return nil, err
	}
	deviceID, err := data.AsInt(did)
	if err != nil {
		return nil, err
	}

	w, err := params.Get(widthPath)
	if err != nil {
		w = data.Int(0) // will be ignored
	}
	width, err := data.AsInt(w)
	if err != nil {
		return nil, err
	}

	h, err := params.Get(heightPath)
	if err != nil {
		h = data.Int(0) // will be ignored
	}
	height, err := data.AsInt(h)
	if err != nil {
		return nil, err
	}

	f, err := params.Get(fpsPath)
	if err != nil {
		f = data.Int(0) // will be ignored
	}
	fps, err := data.AsInt(f)
	if err != nil {
		return nil, err
	}

	cid, err := params.Get(cameraIDPath)
	if err != nil {
		cid = data.Int(0)
	}
	cameraID, err := data.AsInt(cid)
	if err != nil {
		return nil, err
	}

	cs := &captureFromDevice{
		deviceID: deviceID,
		width:    width,
		height:   height,
		fps:      fps,
		cameraID: cameraID,
	}
	// Use ImplementSourceStop helper that can enable this source to stop
	// thread-safe.
	return core.ImplementSourceStop(cs), nil
}

type captureFromDevice struct {
	rwm sync.RWMutex

	deviceID int64
	width    int64
	height   int64
	fps      int64
	cameraID int64
}

// GenerateStream streams video capture datum. OpenCV parameters
// (e.g width, height...) are set when the source is initialized.
//
// Output:
//  capture:   The frame image binary data ('data.Blob'), serialized from
//             OpenCV's matrix data format (`cv::Mat_<cv::Vec3b>`).
//  camera_id: The camera ID.
//  timestamp: The timestamp of capturing. (reed below details)
func (c *captureFromDevice) GenerateStream(ctx *core.Context, w core.Writer) error {
	grabQuit := make(chan struct{}, 1)
	stopped := make(chan struct{}, 1)
	vcap := bridge.NewVideoCapture()
	rootBuf := bridge.NewMatVec3b()
	defer func() {
		grabQuit <- struct{}{}
		<-stopped
		vcap.Delete()
		rootBuf.Delete()
	}()

	if ok := vcap.OpenDevice(int(c.deviceID)); !ok {
		stopped <- struct{}{}
		return fmt.Errorf("error opening device: %v", c.deviceID)
	}

	// OpenCV video capture configuration
	if c.width > 0 {
		vcap.Set(bridge.CvCapPropFrameWidth, int(c.width))
	}
	if c.height > 0 {
		vcap.Set(bridge.CvCapPropFrameHeight, int(c.height))
	}
	if c.fps > 0 {
		vcap.Set(bridge.CvCapPropFps, int(c.fps))
	}

	// read camera frames
	type ret struct {
		buf *bridge.MatVec3b
		err error
	}
	ch := make(chan *ret, 10)
	go func(buf *bridge.MatVec3b) {
		for {
			err := c.grab(&vcap, buf)
			select {
			case <-grabQuit:
				stopped <- struct{}{}
				return
			case ch <- &ret{buf, err}:
			}
			if err != nil {
				stopped <- struct{}{}
				return
			}
		}
	}(&rootBuf)

	// streaming, capture from rootBuf
	buf := bridge.NewMatVec3b()
	defer buf.Delete()
	ctx.Log().Infof("start reading camera device: %v", c.deviceID)
	for {
		res := <-ch
		if res.err != nil {
			return res.err
		}
		func() {
			c.rwm.RLock()
			defer c.rwm.RUnlock()
			res.buf.CopyTo(&buf)
		}()
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
		if err := w.Write(ctx, &t); err != nil {
			return err
		}
	}
	return nil
}

func (c *captureFromDevice) grab(vcap *bridge.VideoCapture, buf *bridge.MatVec3b) error {
	if !vcap.IsOpened() {
		return fmt.Errorf("video stream or file closed, device no: %d)",
			c.deviceID)
	}
	tmpBuf := bridge.NewMatVec3b()
	defer tmpBuf.Delete()
	if ok := vcap.Read(tmpBuf); !ok {
		return fmt.Errorf("cannot read a new file (device no: %d)", c.deviceID)
	}

	c.rwm.Lock()
	defer c.rwm.Unlock()
	tmpBuf.CopyTo(buf)
	return nil
}

func (c *captureFromDevice) Stop(ctx *core.Context) error {
	return nil
}

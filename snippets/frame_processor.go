package snippets

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type FrameApplier struct {
	config string
	fp     bridge.FrameProcessor
}

type CaptureInfo struct {
	capture  []byte
	cameraID int64
}

func (fa *FrameApplier) Init(ctx *core.Context) error {
	fa.fp = bridge.NewFrameProcessor(fa.config)
	return nil
}

func (fa *FrameApplier) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	info, err := getCapture(t)
	if err != nil {
		return err
	}
	now := time.Now()
	inow := now.UnixNano() / int64(time.Millisecond) // [ms]
	f := fa.fp.Apply(bridge.DeserializeMatVec3b(info.capture), inow, int(info.cameraID))
	t.Data["frame"] = tuple.Blob(f.Serialize())

	return nil
}

func getCapture(t *tuple.Tuple) (CaptureInfo, error) {
	info := CaptureInfo{}
	c, err := t.Data.Get("capture")
	if err != nil {
		return info, fmt.Errorf("cannot get capture data")
	}
	capture, err := tuple.AsBlob(c)
	if err != nil {
		return info, fmt.Errorf("capture data must be byte array type")
	}

	id, err := t.Data.Get("camera_id")
	if err != nil {
		return info, fmt.Errorf("cannot get camera id")
	}
	cameraID, err := tuple.AsInt(id)
	if err != nil {
		return info, fmt.Errorf("camera ID must be int type")
	}
	return CaptureInfo{
		capture:  capture,
		cameraID: cameraID,
	}, nil
}

func (fa *FrameApplier) Terminate(ctx *core.Context) error {
	fa.fp.Delete()
	return nil
}

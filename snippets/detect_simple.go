package snippets

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type DetectSimple struct {
	ConfigPath string
	Config     *conf.DetectSimpleConfig
	detector   *bridge.Detector
	lastFrame  map[int64]*tuple.Tuple
}

func (d *DetectSimple) Init(ctx *core.Context) error {
	detectConfig, err := conf.GetDetectSimpleSnippetConfig(d.ConfigPath)
	if err != nil {
		return err
	}
	d.Config = &detectConfig
	detector := bridge.NewDetector(detectConfig.DetectorConfig)
	d.detector = &detector
	return nil
}

func (d *DetectSimple) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	switch t.InputName {
	case "frame":
		cameraId, err := t.Data.Get("camera_id")
		if err != nil {
			return err
		}
		id, err := cameraId.AsInt()
		if err != nil {
			return err
		}
		d.lastFrame[id] = t
	case "tick":
		if len(d.lastFrame) == 0 {
			return nil
		}
		for _, fTuple := range d.lastFrame {
			f, err := getFrame(fTuple)
			if err != nil {
				return err
			}

			fPointer := bridge.DeserializeFrame(f)
			defer fPointer.Delete()
			s := time.Now().UnixNano() / int64(time.Millisecond)
			drPointer := d.detector.Detect(fPointer)

			fTuple.Data["detection_result"] = tuple.Blob(drPointer.Serialize())
			fTuple.Data["detection_time"] = tuple.Timestamp(fTuple.Timestamp) // same as frame create time

			if d.Config.PlayerFlag {
				ms := time.Now().UnixNano()/int64(time.Millisecond) - s
				drw := bridge.DetectDrawResult(fPointer, drPointer, ms)
				fTuple.Data["detection_draw_result"] = tuple.Blob(drw.ToJpegData(d.Config.JpegQuality))
			}

			w.Write(ctx, fTuple)
		}
	}
	return nil
}

func getFrame(t *tuple.Tuple) ([]byte, error) {
	f, err := t.Data.Get("frame")
	if err != nil {
		return []byte{}, fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return []byte{}, fmt.Errorf("frame data must be byte array type")
	}
	return frame, nil
}

func (d *DetectSimple) Terminate(ctx *core.Context) error {
	d.detector.Delete()
	return nil
}

package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type RecognizeCaffeConfig struct {
	PlayerFlag bool
}

type RecognizeCaffe struct {
	Config        RecognizeCaffeConfig
	frameInfoChan chan FrameInfo
}

type FrameInfo struct {
	index int
	fr    bridge.Frame
	dr    bridge.DetectionResult
}

func (rc *RecognizeCaffe) Init(ctx *core.Context) error {
	return nil
}

func (rc *RecognizeCaffe) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	f, err := t.Data.Get("frame")
	if err != nil {
		return fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return fmt.Errorf("frame data must be byte array type")
	}

	d, err := t.Data.Get("detection_result")
	if err != nil {
		return fmt.Errorf("cannot get detection result")
	}
	detectionResult, err := d.AsBlob()
	if err != nil {
		return fmt.Errorf("detection result data must be byte array type")
	}

	fr := bridge.ConvertToFramePointer(frame)
	dr := bridge.ConvertToDetectionResultPointer(detectionResult)

	governor(fr, dr, rc)
	recognize(fr, dr, t, rc)

	w.Write(ctx, t)
	return nil
}

func governor(fr bridge.Frame, dr bridge.DetectionResult, rc *RecognizeCaffe) {
	// join where meta.time is equal
}

func recognize(fr bridge.Frame, dr bridge.DetectionResult, t *tuple.Tuple, rc *RecognizeCaffe) {
	var taggers bridge.ImageTaggerCaffes
	bridge.ImageTaggerCaffe_SetUp(taggers, nil) // TODO set up recognize configuration
	recogDr, recogDrByte := bridge.ImageTaggerCaffe_PredictTagsBatch(taggers, fr, dr)

	t.Data["recognize_detection_result"] = tuple.Blob(recogDrByte)

	if rc.Config.PlayerFlag {
		drwResult := bridge.RecognizeDrawResult(fr, recogDr)
		t.Data["recognize_draw_result"] = tuple.Blob(drwResult)
	}
}

func (rc *RecognizeCaffe) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (rc *RecognizeCaffe) OutputSchema(ss []*core.Schema) (*core.Schema, error) {
	return nil, nil
}

func (rc *RecognizeCaffe) Terminate(ctx *core.Context) error {
	return nil
}

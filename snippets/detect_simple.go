package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type DetectSimpleConfig struct {
	PlayerFlag bool
}

type DetectSimple struct {
	Config   DetectSimpleConfig
	detector bridge.Detector
}

func (d *DetectSimple) Init(ctx *core.Context) error {
	var detector bridge.Detector
	bridge.Detector_SetUp(detector, nil) // TODO setup detector config
	d.detector = detector
	return nil
}

func (d *DetectSimple) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	f, err := t.Data.Get("frame")
	if err != nil {
		return fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return fmt.Errorf("frame data must be byte array type")
	}

	frPointer := bridge.ConvertToFramePointer(frame)
	s := bridge.Scouter_GetEpochms()
	dr, drByte := bridge.Detector_Detect(d.detector, frPointer)

	t.Data["detection_result"] = tuple.Blob(drByte)
	t.Data["detection_time"] = tuple.Timestamp(t.Timestamp) // same as frame create time

	if d.Config.PlayerFlag {
		ms := bridge.Scouter_GetEpochms() - s
		drwByte := bridge.DetectDrawResult(frPointer, dr, ms)
		t.Data["detection_draw_result"] = tuple.Blob(drwByte)
	}

	w.Write(ctx, t)
	return nil
}

func (d *DetectSimple) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (d *DetectSimple) OutputSchema(ss []*core.Schema) (*core.Schema, error) {
	return nil, nil
}

package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type DetectSimpleConfig struct {
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
	m, err := t.Data.AsMap()
	if err != nil {
		return fmt.Errorf("cannot get frame data")
	}
	f, err := m.Get("frame")
	if err != nil {
		return fmt.Errorf("cannot get frame data in map")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return fmt.Errorf("frame data must be []byte type")
	}

	frPointer := bridge.ConvertToFramePointer(frame)
	s := bridge.Scouter_GetEpochms()
	dr, derByte := bridge.Detector_Detect(d.detector, frPointer)
	ms := bridge.Scouter_GetEpochms() - s
	_, drwByte := bridge.DetectDrawResult(frPointer, dr, ms)

	m["detection_result"] = tuple.Blob(drByte)
	m["result_frame"] = tuple.Blob(drwByte)

	return nil
}

func (d *DetectSimple) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (d *DetectSimple) OutputSchema(ss []*core.Schema) (*core.Schema, error) {
	return nil, nil
}

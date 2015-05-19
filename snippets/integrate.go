package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type IntegrateConfig struct {
}

type Integrate struct {
	Config IntegrateConfig
}

func (itr *Integrate) Init(ctx *core.Context) error {
	return nil
}

func (itr *Integrate) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	f, err := t.Data.Get("frame")
	if err != nil {
		return fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return fmt.Errorf("frame data must be byte array type")
	}

	d, err := t.Data.Get("recognize_detection_result")
	if err != nil {
		return fmt.Errorf("cannot get detection result")
	}
	detectionResult, err := d.AsBlob()
	if err != nil {
		return fmt.Errorf("detection result data must be byte array type")
	}

	fr := bridge.ConvertToFramePointer(frame)
	dr := bridge.ConvertToDetectionResultPointer(detectionResult)

	tracking(fr, dr, itr)

	return nil
}

func tracking(fr bridge.Frame, dr bridge.DetectionResult, itr *Integrate) {

}

func (itr *Integrate) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (itr *Integrate) OutputSchema([]*core.Schema) (*core.Schema, error) {
	return nil, nil
}

package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type IntegrateConfig struct {
	PlayerFlag bool
}

type Integrate struct {
	Config     IntegrateConfig
	integrator bridge.Integrator
}

func (itr *Integrate) Init(ctx *core.Context) error {
	var integrator bridge.Integrator
	bridge.IntegratorSetUp(integrator, nil) // TODO configuration
	itr.integrator = integrator
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

	bridge.Integrator_Push(itr.integrator, fr, dr)
	if bridge.Integrator_TrackerReady(itr.integrator) {
		return nil // TODO set empty tracking result?
	}

	_, trByte := bridge.Integrator_Track(itr.integrator)
	t.Data["tracking_result"] = tuple.Blob(trByte)

	if itr.Config.PlayerFlag {
		// TODO draw result for debug
	}

	w.Write(ctx, t)
	return nil
}

func (itr *Integrate) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (itr *Integrate) OutputSchema([]*core.Schema) (*core.Schema, error) {
	return nil, nil
}

func (itr *Integrate) Terminate(ctx *core.Context) error {
	return nil
}

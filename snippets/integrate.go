package snippets

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type Integrate struct {
	ConfigPath string
	Config     conf.IntegrateConfig
	integrator bridge.Integrator
}

type TrackingInfo struct {
	index int
	fr    []byte
	dr    []byte
}

func (itr *Integrate) Init(ctx *core.Context) error {
	config, err := conf.GetIntegrateConfig(itr.ConfigPath)
	if err != nil {
		return err
	}
	itr.Config = config
	integrator := bridge.NewIntegrator(config.IntegrateConfig)
	itr.integrator = integrator
	return nil
}

func (itr *Integrate) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	fi, err := getTrackingInfo(t)
	if err != nil {
		return nil
	}

	fr := bridge.DeserializeFrame(fi.fr)
	dr := bridge.DeserializeDetectionResult(fi.dr)

	itr.integrator.Integrator_Push(fr, dr)
	if !itr.integrator.Integrator_TrackerReady() {
		return nil // TODO set empty tracking result?
	}

	tr := itr.integrator.Integrator_Track()
	defer tr.Delete()
	t.Data["tracking_result"] = tuple.Blob(tr.Serialize())

	if itr.Config.PlayerFlag {
		// TODO draw result for debug
	}

	w.Write(ctx, t)
	return nil
}

func getTrackingInfo(t *tuple.Tuple) (TrackingInfo, error) {
	f, err := t.Data.Get("frame")
	if err != nil {
		return TrackingInfo{}, fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return TrackingInfo{}, fmt.Errorf("frame data must be byte array type")
	}

	d, err := t.Data.Get("detection_result")
	if err != nil {
		return TrackingInfo{}, fmt.Errorf("cannot get detection result")
	}
	detectionResult, err := d.AsBlob()
	if err != nil {
		return TrackingInfo{}, fmt.Errorf("detection result data must be byte array type")
	}

	return TrackingInfo{
		fr: frame,
		dr: detectionResult}, nil
}

func (itr *Integrate) Terminate(ctx *core.Context) error {
	itr.integrator.Delete()
	return nil
}

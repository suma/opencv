package snippets

import (
	"fmt"
	"io/ioutil"
	"os"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type Integrate struct {
	ConfigPath      string
	Config          conf.IntegrateConfig
	integrator      bridge.Integrator
	instanceManager bridge.InstanceManager
	visualizer      bridge.Visualizer
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
	itr.integrator = bridge.NewIntegrator(config.IntegratorConfig)
	itr.instanceManager = bridge.NewInstanceManager(config.InstanceManagerConfig)
	itr.visualizer = bridge.NewVisualizer(config.VisualizerConfig, itr.instanceManager)
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
	currentStates := itr.instanceManager.GetCurrentStates(tr)
	defer currentStates.Delete()

	now := t.Timestamp.UnixNano() / int64(time.Millisecond)
	statesJson := currentStates.ConvertSatesToJson(itr.Config.FloorID, now)

	t.Data["instance_states"] = tuple.String(statesJson)

	if itr.Config.PlayerFlag {
		trajectories := itr.visualizer.PlotTrajectories()
		for i, traj := range trajectories {
			key := fmt.Sprintf("integrate_result[%d]", i)
			t.Data[key] = tuple.Blob(traj.ToJpegData(itr.Config.JpegQuality))
			// following is debug for scouter integrator
			s := time.Now().UnixNano() / int64(time.Millisecond)
			ioutil.WriteFile(fmt.Sprintf("./integrate_%d_%v.jpg", i, fmt.Sprint(s)),
				traj.ToJpegData(50), os.ModePerm)
		}
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
	itr.visualizer.Delete()
	itr.instanceManager.Delete()
	itr.integrator.Delete()
	return nil
}

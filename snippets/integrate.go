package snippets

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

// Integrate is several detected frames.
type Integrate struct {
	// ConfigPath is the path of external configuration file.
	ConfigPath string

	config          conf.IntegrateConfig
	integrator      bridge.Integrator
	instanceManager bridge.InstanceManager
	visualizer      bridge.Visualizer
}

// trackingInfo is pair of frame and detection result data.
type trackingInfo struct {
	index int
	fr    []byte
	dr    []byte
}

// Init prepares integration information set by external configuration file.
func (itr *Integrate) Init(ctx *core.Context) error {
	config, err := conf.GetIntegrateConfig(itr.ConfigPath)
	if err != nil {
		return err
	}
	itr.config = config
	itr.integrator = bridge.NewIntegrator(config.IntegratorConfig)
	itr.instanceManager = bridge.NewInstanceManager(config.InstanceManagerConfig)
	itr.visualizer = bridge.NewVisualizer(config.VisualizerConfig, itr.instanceManager)
	return nil
}

// Process add integration information to frames. Integration is caching several
// frame data.
func (itr *Integrate) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	fi, err := getTrackingInfo(t)
	if err != nil {
		return nil
	}

	fr := bridge.DeserializeFrame(fi.fr)
	dr := bridge.DeserializeDetectionResult(fi.dr)

	itr.integrator.Integrator_Push(fr, dr)
	if !itr.integrator.Integrator_TrackerReady() {
		w.Write(ctx, t)
		return nil
	}

	tr := itr.integrator.Integrator_Track()
	defer tr.Delete()
	currentStates := itr.instanceManager.GetCurrentStates(tr)
	defer currentStates.Delete()

	now := t.Timestamp.UnixNano() / int64(time.Millisecond)
	statesJSON := currentStates.ConvertSatesToJson(itr.config.FloorID, now)

	t.Data["instance_states"] = tuple.String(statesJSON)

	if itr.config.PlayerFlag {
		trajectories := itr.visualizer.PlotTrajectories()
		debugArray := tuple.Array([]tuple.Value{})
		for _, traj := range trajectories {
			jpeg := tuple.Blob(traj.ToJpegData(itr.config.JpegQuality))
			debugArray = append(debugArray, jpeg)
		}
		t.Data["integrate_result"] = debugArray
	}

	w.Write(ctx, t)
	return nil
}

func getTrackingInfo(t *tuple.Tuple) (trackingInfo, error) {
	f, err := t.Data.Get("frame")
	if err != nil {
		return trackingInfo{}, fmt.Errorf("cannot get frame data")
	}
	frame, err := tuple.AsBlob(f)
	if err != nil {
		return trackingInfo{}, fmt.Errorf("frame data must be byte array type")
	}

	d, err := t.Data.Get("detection_result")
	if err != nil {
		return trackingInfo{}, fmt.Errorf("cannot get detection result")
	}
	detectionResult, err := tuple.AsBlob(d)
	if err != nil {
		return trackingInfo{}, fmt.Errorf("detection result data must be byte array type")
	}

	return trackingInfo{
		fr: frame,
		dr: detectionResult}, nil
}

// Terminate this component.
func (itr *Integrate) Terminate(ctx *core.Context) error {
	itr.visualizer.Delete()
	itr.instanceManager.Delete()
	itr.integrator.Delete()
	return nil
}

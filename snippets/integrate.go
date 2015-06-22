package snippets

import (
	"fmt"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"sync"
	"time"
)

// Integrate is several detected frames.
type Integrate struct {
	// ConfigPath is the path of external configuration file.
	ConfigPath string

	config            conf.IntegrateConfig
	integrator        bridge.Integrator
	instanceManager   bridge.InstanceManager
	visualizer        bridge.Visualizer
	trackingInfoQueue map[string]map[time.Time]trackingInfo
	mu                sync.RWMutex
}

// trackingInfo is pair of frame and detection result data.
type trackingInfo struct {
	name       string
	fr         []byte
	dr         []byte
	detectTime time.Time
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
	itr.trackingInfoQueue = map[string]map[time.Time]trackingInfo{}
	itr.mu = sync.RWMutex{}
	return nil
}

func (itr *Integrate) aggregation(t trackingInfo) {
	if len(itr.config.FrameInputKeys) == 1 {
		return
	}

	queue, ok := itr.trackingInfoQueue[t.name]
	if !ok {
		queue = map[time.Time]trackingInfo{}
	}
	queue[t.detectTime] = t
	itr.trackingInfoQueue[t.name] = queue
}

func (itr *Integrate) pop(t trackingInfo) (bool, []trackingInfo) {
	inputSource := len(itr.config.FrameInputKeys)
	if inputSource == 1 {
		return true, []trackingInfo{t}
	}

	trackingInfos := []trackingInfo{}
	for _, v := range itr.trackingInfoQueue {
		ti, ok := v[t.detectTime]
		if ok {
			trackingInfos = append(trackingInfos, ti)
		}
	}
	if len(trackingInfos) == inputSource {
		return true, trackingInfos
	}
	return false, []trackingInfo{}
}

// Process add integration information to frames. Integration is caching several
// frame data.
func (itr *Integrate) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	fi, err := getTrackingInfo(t)
	if err != nil {
		return nil
	}
	itr.aggregation(fi)
	ok, infos := itr.pop(fi)
	if !ok {
		return nil
	}

	frs := []bridge.Frame{}
	drs := []bridge.DetectionResult{}
	for _, ti := range infos {
		frs = append(frs, bridge.DeserializeFrame(ti.fr))
		drs = append(drs, bridge.DeserializeDetectionResult(ti.dr))
	}

	itr.integrator.Integrator_Push(frs, drs)
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
	info := trackingInfo{}
	f, err := t.Data.Get("frame")
	if err != nil {
		return info, fmt.Errorf("cannot get frame data")
	}
	frame, err := tuple.AsBlob(f)
	if err != nil {
		return info, fmt.Errorf("frame data must be byte array type")
	}

	d, err := t.Data.Get("detection_result")
	if err != nil {
		return info, fmt.Errorf("cannot get detection result")
	}
	detectionResult, err := tuple.AsBlob(d)
	if err != nil {
		return info, fmt.Errorf("detection result data must be byte array type")
	}

	ti, err := t.Data.Get("detection_time")
	if err != nil {
		return info, fmt.Errorf("cannot get frame detection time")
	}
	detectTime, err := tuple.AsTimestamp(ti)
	if err != nil {
		return info, fmt.Errorf("detection time must be timestamp type")
	}

	return trackingInfo{
		name:       t.InputName,
		fr:         frame,
		dr:         detectionResult,
		detectTime: detectTime}, nil
}

// Terminate this component.
func (itr *Integrate) Terminate(ctx *core.Context) error {
	itr.visualizer.Delete()
	itr.instanceManager.Delete()
	itr.integrator.Delete()
	return nil
}

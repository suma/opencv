package integrator

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type framesTrackerUDSF struct {
	tracker                   bridge.Tracker
	instanceManager           bridge.InstanceManager
	instanceStatesIDFieldName string
	framesFieldName           string
	cameraIDFieldName         string
	imageFieldName            string
	mvRegionsFieldName        string
	timestampFieldName        string
}

func (sf *framesTrackerUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {

	// instance id
	isID, err := t.Data.Get(sf.instanceStatesIDFieldName)
	if err != nil {
		return err
	}

	// timestamp
	ts, err := t.Data.Get(sf.timestampFieldName)
	if err != nil {
		return err
	}
	frameTime, err := data.AsTimestamp(ts)
	if err != nil {
		return err
	}

	// multi place frames
	frames, err := t.Data.Get(sf.framesFieldName)
	if err != nil {
		return err
	}
	frameArray, err := data.AsArray(frames)
	if err != nil {
		return err
	}

	matMap, err := sf.convertToMatVecMap(frameArray)
	defer func() {
		for _, v := range matMap {
			v.Delete()
		}
	}()
	if err != nil {
		return err
	}

	// moving detection result
	mvRegions, err := t.Data.Get(sf.mvRegionsFieldName)
	if err != nil {
		return err
	}
	mvRegionsArray, err := data.AsArray(mvRegions)
	if err != nil {
		return err
	}

	mvCans, err := convertToMVCandidateSlice(mvRegionsArray)
	defer func() {
		for _, c := range mvCans {
			c.Delete()
		}
	}()
	if err != nil {
		return err
	}

	timestamp := time.Duration(frameTime.UnixNano()) / time.Millisecond
	sf.tracker.Push(matMap, mvCans, uint64(timestamp))

	if sf.tracker.Ready() {
		tr := sf.tracker.Track(uint64(timestamp))
		sf.instanceManager.Update(tr)

		currentStates := sf.instanceManager.GetCurrentStates()
		if len(currentStates) <= 0 {
			ctx.Log().Info("tracking current status is empty")
			return nil
		}
		defer func() {
			for _, s := range currentStates {
				s.Delete()
			}
		}()
		ctx.Log().Infof("tracking is completed:%d", len(currentStates))

		traceCopyFlag := len(t.Trace) > 0
		for _, s := range currentStates {
			now := time.Now()
			m := data.Map{
				"states_id":      isID,
				"states_count":   data.Int(len(currentStates)),
				"instance_state": data.Blob(s.Serialize()),
				"timestamp":      data.Timestamp(frameTime),
			}
			traces := []core.TraceEvent{}
			if traceCopyFlag { // reduce copy cost when trace mode is off
				traces = make([]core.TraceEvent, len(t.Trace), (cap(t.Trace)+1)*2)
				copy(traces, t.Trace)
			}
			tu := &core.Tuple{
				Data:          m,
				Timestamp:     now,
				ProcTimestamp: t.ProcTimestamp,
				Trace:         traces,
			}
			w.Write(ctx, tu)
		}
	}
	return nil
}

func (sf *framesTrackerUDSF) convertToMatVecMap(frameArray data.Array) (
	map[int]bridge.MatVec3b, error) {

	matMap := map[int]bridge.MatVec3b{}
	for _, f := range frameArray {
		fMap, err := data.AsMap(f)
		if err != nil {
			return nil, err
		}

		id, err := fMap.Get(sf.cameraIDFieldName)
		if err != nil {
			return nil, err
		}
		cameraID, err := data.AsInt(id)
		if err != nil {
			return nil, err
		}

		image, err := fMap.Get(sf.imageFieldName)
		if err != nil {
			return nil, err
		}
		imageByte, err := data.AsBlob(image)
		if err != nil {
			return nil, err
		}

		matMap[int(cameraID)] = bridge.DeserializeMatVec3b(imageByte)
	}
	return matMap, nil
}

func convertToMVCandidateSlice(mvRegionsArray data.Array) (
	[]bridge.MVCandidate, error) {

	mvCans := []bridge.MVCandidate{}
	for _, r := range mvRegionsArray {
		b, err := data.AsBlob(r)
		if err != nil {
			return nil, err
		}
		mvCans = append(mvCans, bridge.DeserializeMVCandiate(b))
	}
	return mvCans, nil
}

func (sf *framesTrackerUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createFramesTrackerUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	trackerParam string, instanceManagerParam string, stream string,
	instanceStatesIDFieldName string, framesFieldName string,
	cameraIDFieldName string, imageFieldname string, mvRegionsFieldName string,
	timestampFieldName string) (
	udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "frame_tracker_stream",
	}); err != nil {
		return nil, err
	}

	trackerState, err := lookupTrackerParamState(ctx, trackerParam)
	if err != nil {
		return nil, err
	}

	instanceManagerState, err := lookupInstanceManagerParamState(
		ctx, instanceManagerParam)
	if err != nil {
		return nil, err
	}

	return &framesTrackerUDSF{
		tracker:                   trackerState.t,
		instanceManager:           instanceManagerState.m,
		instanceStatesIDFieldName: instanceStatesIDFieldName,
		framesFieldName:           framesFieldName,
		cameraIDFieldName:         cameraIDFieldName,
		imageFieldName:            imageFieldname,
		mvRegionsFieldName:        mvRegionsFieldName,
		timestampFieldName:        timestampFieldName,
	}, nil
}

// FramesTrackerStreamFuncCreator is a creator of frame tracking UDSF.
type FramesTrackerStreamFuncCreator struct{}

// CreateStreamFunction creates instance state from tracked detections.
// This function need moving matched detection datum per captured frame.
// If captured frames include multiple places, then frames and detections could
// be distinguished with camera ID.
//
// Input tuples are required to have following `data.Map` structure, each key
// name is addressed with UDSF's arguments.
//
//  data.Map{
//    "instanceStatesIDFieldName": [ID],
//    "framesFieldName"          : data.Array{
//      []data.Map{
//        "cameraIDFieldName": [camera ID],
//        "imageFiledname"   : [image data] (data.Blob),
//      }
//    },
//    "mvRegionsFieldName": [moving matched detection result] ([]data.Blob),
//    "timestampFieldName": [frame captured time] (data.Timestamp)
//  }
func (c *FramesTrackerStreamFuncCreator) CreateStreamFunction() interface{} {
	return createFramesTrackerUDSF
}

func (c *FramesTrackerStreamFuncCreator) TypeName() string {
	return "multi_region_tracking"
}

func lookupTrackerParamState(ctx *core.Context, trackerParam string) (
	*TrackerParamState, error) {
	st, err := ctx.SharedStates.Get(trackerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*TrackerParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to tracker_param.state", trackerParam)
}

func lookupInstanceManagerParamState(ctx *core.Context, instanceManagerParam string) (
	*InstanceManagerParamState, error) {
	st, err := ctx.SharedStates.Get(instanceManagerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*InstanceManagerParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to instance_manager_param.state",
		instanceManagerParam)
}

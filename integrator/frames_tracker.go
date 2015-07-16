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
	tracker            bridge.Tracker
	instanceManager    bridge.InstanceManager
	framesFieldName    string
	cameraIDFieldName  string
	imageFieldName     string
	mvRegionsFieldName string
}

func (sf *framesTrackerUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	frames, err := t.Data.Get(sf.framesFieldName)
	if err != nil {
		return err
	}
	frameArray, err := data.AsArray(frames)
	if err != nil {
		return err
	}

	matMap, err := sf.convertToMatVecMap(frameArray)
	if err != nil {
		return err
	}

	mvRegions, err := t.Data.Get(sf.mvRegionsFieldName)
	if err != nil {
		return err
	}
	mvRegionsArray, err := data.AsArray(mvRegions)
	if err != nil {
		return err
	}

	mvCans, err := convertToMVCandidateSlice(mvRegionsArray)
	if err != nil {
		return err
	}

	timestamp := time.Duration(t.Timestamp.UnixNano()) / time.Microsecond
	sf.tracker.Push(matMap, mvCans, uint64(timestamp))

	if sf.tracker.Ready() {
		tr := sf.tracker.Track(uint64(timestamp))
		sf.instanceManager.Updaate(tr)

		currentStates := sf.instanceManager.GetCurrentStates()
		for _, s := range currentStates {
			now := time.Now()
			m := data.Map{
				"instance_state": data.Blob(s.Serialize()),
			}
			tu := &core.Tuple{
				Data:          m,
				Timestamp:     now,
				ProcTimestamp: t.ProcTimestamp,
				Trace:         make([]core.TraceEvent, 0),
			}
			w.Write(ctx, tu)
		}
	}
	return nil
}

func (sf *framesTrackerUDSF) convertToMatVecMap(frameArray data.Array) (map[int]bridge.MatVec3b, error) {
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

func convertToMVCandidateSlice(mvRegionsArray data.Array) ([]bridge.MVCandidate, error) {
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

func createFramesTrackerUDSF(ctx *core.Context, decl udf.UDSFDeclarer, trackerParam string,
	instanceManagerParam string, stream string, framesFieldName string,
	cameraIDFieldName string, imageFieldname string, mvRegionsFieldName string) (
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

	instanceManagerState, err := lookupInstanceManagerParamState(ctx, instanceManagerParam)
	if err != nil {
		return nil, err
	}

	return &framesTrackerUDSF{
		tracker:            trackerState.t,
		instanceManager:    instanceManagerState.m,
		framesFieldName:    framesFieldName,
		cameraIDFieldName:  cameraIDFieldName,
		imageFieldName:     imageFieldname,
		mvRegionsFieldName: mvRegionsFieldName,
	}, nil
}

type FramesTrackerStreamFuncCreator struct{}

func (c *FramesTrackerStreamFuncCreator) CreateStreamFunction() interface{} {
	return createFramesTrackerUDSF
}

func (c *FramesTrackerStreamFuncCreator) TypeName() string {
	return "tracking"
}

func lookupTrackerParamState(ctx *core.Context, trackerParam string) (*TrackerParamState, error) {
	st, err := ctx.SharedStates.Get(trackerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*TrackerParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to tracker_param.state", trackerParam)
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
	return nil, fmt.Errorf("state '%v' cannot be converted to instance_manager_param.state",
		instanceManagerParam)
}

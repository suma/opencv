package integrator

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type TrackInstanceStatesUDSFCreator struct{}

//  data.Map{
//    "instanceStatesIDFieldName": [ID],
//    "framesFieldName"          : data.Array{
//      []data.Map{
//        "cameraIDFieldName": [camera ID],
//        "imageFiledname"   : [image data] (data.Blob),
//      }
//    },
//    "trackeesFieldName" : [tracking result] ([]data.Blob),
//    "timestampFieldName": [frame captured time[us]] (data.Int)
//  }
func (c *TrackInstanceStatesUDSFCreator) CreateStreamFunction() interface{} {
	return createTrackInstanceStatesUDSF
}

func (c *TrackInstanceStatesUDSFCreator) TypeName() string {
	return "tracking_instance_states"
}

func createTrackInstanceStatesUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	instanceManagerParam string, stream string,
	instanceStatesIDFieldName string, framesFieldName string,
	cameraIDFieldName string, imageFieldName string, trackeesFieldName string,
	timestampFieldName string) (udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "track_instance_states_stream",
	}); err != nil {
		return nil, err
	}

	imState, err := lookupInstanceManagerParamState(ctx, instanceManagerParam)
	if err != nil {
		return nil, err
	}

	return &trackInstanceStatesUDSF{
		instanceManager:           &imState.m,
		instanceStatesIDFieldName: instanceStatesIDFieldName,
		framesFieldName:           framesFieldName,
		cameraIDFieldName:         cameraIDFieldName,
		imageFieldName:            imageFieldName,
		trackeesFieldName:         trackeesFieldName,
		timestampFieldName:        timestampFieldName,
	}, nil
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

type trackInstanceStatesUDSF struct {
	instanceManager           *bridge.InstanceManager
	instanceStatesIDFieldName string
	framesFieldName           string
	cameraIDFieldName         string
	imageFieldName            string
	trackeesFieldName         string
	timestampFieldName        string
}

func (sf *trackInstanceStatesUDSF) Process(ctx *core.Context, t *core.Tuple,
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
	trTime, err := data.AsInt(ts)
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

	// trackees
	trs, err := t.Data.Get(sf.trackeesFieldName)
	if err != nil {
		return err
	}
	trsArray, err := data.AsArray(trs)
	if err != nil {
		return err
	}
	trackees, err := convertToTrackeeSlice(trsArray)
	if err != nil {
		return err
	}
	defer func() {
		for _, tr := range trackees {
			tr.MVCandidate.Delete()
		}
	}()

	sf.instanceManager.Update(matMap, trackees, uint64(trTime))

	states := sf.instanceManager.GetCurrentStates()
	if len(states) <= 0 {
		ctx.Log().Info("instance states is empty")
		return nil
	}
	defer func() {
		for _, s := range states {
			s.Delete()
		}
	}()

	traceCopyFlag := len(t.Trace) > 0
	for _, state := range states {
		now := time.Now()
		m := data.Map{
			"states_id":   isID,
			"state_count": data.Int(len(states)),
			"state":       data.Blob(state.Serialize()),
			"timestamp":   ts,
		}
		traces := []core.TraceEvent{}
		if traceCopyFlag {
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
	return nil
}

func (sf *trackInstanceStatesUDSF) convertToMatVecMap(frameArray data.Array) (
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

func convertToTrackeeSlice(trsArray data.Array) ([]bridge.Trackee, error) {
	trackees := []bridge.Trackee{}
	for _, tr := range trsArray {
		trMap, err := data.AsMap(tr)
		if err != nil {
			return nil, err
		}
		var colorID uint64
		if cid, err := trMap.Get("color_id"); err != nil {
			return nil, err
		} else if cidInt, err := data.AsInt(cid); err != nil {
			return nil, err
		} else {
			colorID = uint64(cidInt)
		}
		var mvRegion bridge.MVCandidate
		if mvCan, err := trMap.Get("moving_detected"); err != nil {
			return nil, err
		} else if mvByte, err := data.AsBlob(mvCan); err != nil {
			return nil, err
		} else {
			mvRegion = bridge.DeserializeMVCandiate(mvByte)
		}
		var interpolated bool
		if interpo, err := trMap.Get("interpolated"); err != nil {
			return nil, err
		} else if interpolated, err = data.AsBool(interpo); err != nil {
			return nil, err
		}
		var timestamp uint64
		if ts, err := trMap.Get("timestamp"); err != nil {
			return nil, err
		} else if tsInt, err := data.AsInt(ts); err != nil {
			return nil, err
		} else {
			timestamp = uint64(tsInt)
		}

		trackee := bridge.Trackee{
			ColorID:      colorID,
			MVCandidate:  mvRegion,
			Interpolated: interpolated,
			Timestamp:    timestamp,
		}

		trackees = append(trackees, trackee)
	}
	return trackees, nil
}

func (sf *trackInstanceStatesUDSF) Terminate(ctx *core.Context) error {
	return nil
}

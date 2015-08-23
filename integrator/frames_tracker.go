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
	tracker                   *bridge.Tracker
	instanceStatesIDFieldName data.Path
	framesFieldName           data.Path
	cameraIDFieldName         data.Path
	imageFieldName            data.Path
	mvRegionsFieldName        data.Path
	timestampFieldName        data.Path
}

func (sf *framesTrackerUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {

	// instance id
	isID, err := t.Data.Get(sf.instanceStatesIDFieldName)
	if err != nil {
		return err
	}

	// timestamp
	// ts, err := t.Data.Get(sf.timestampFieldName)
	// if err != nil {
	// 	return err
	// }
	// frameTime, err := data.AsTimestamp(ts)
	// if err != nil {
	// 	return err
	// }

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

	// TODO delete
	// timestamp := time.Duration(frameTime.UnixNano()) / time.Microsecond
	// sf.tracker.Push(matMap, mvCans, uint64(timestamp))

	if sf.tracker.Ready() {
		trs := sf.tracker.Track()
		defer func() {
			for _, tr := range trs {
				tr.MVCandidate.Delete()
			}
		}()

		traceCopyFlag := len(t.Trace) > 0
		for _, trackee := range trs {
			now := time.Now()
			m := data.Map{
				"states_id":       isID,
				"trackee_count":   data.Int(len(trs)),
				"color_id":        data.Int(trackee.ColorID),
				"moving_detected": data.Blob(trackee.MVCandidate.Serialize()),
				"interpolated":    data.Bool(trackee.Interpolated),
				"timestamp":       data.Int(trackee.Timestamp),
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
		mvCans = append(mvCans, bridge.DeserializeMVCandidate(b))
	}
	return mvCans, nil
}

// Terminate the components.
func (sf *framesTrackerUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createFramesTrackerUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	trackerParam string, stream string,
	instanceStatesIDFieldName string, framesFieldName string,
	cameraIDFieldName string, imageFieldname string, mvRegionsFieldName string,
	timestampFieldName string) (udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "frame_tracker_stream",
	}); err != nil {
		return nil, err
	}

	trackerState, err := lookupTrackerParamState(ctx, trackerParam)
	if err != nil {
		return nil, err
	}

	return &framesTrackerUDSF{
		tracker:                   &trackerState.t,
		instanceStatesIDFieldName: data.MustCompilePath(instanceStatesIDFieldName),
		framesFieldName:           data.MustCompilePath(framesFieldName),
		cameraIDFieldName:         data.MustCompilePath(cameraIDFieldName),
		imageFieldName:            data.MustCompilePath(imageFieldname),
		mvRegionsFieldName:        data.MustCompilePath(mvRegionsFieldName),
		timestampFieldName:        data.MustCompilePath(timestampFieldName),
	}, nil
}

// FramesTrackerStreamFuncCreator is a creator of frame tracking UDSF.
type FramesTrackerStreamFuncCreator struct{}

// CreateStreamFunction creates instance state from tracked detections.
// This function need moving matched detection datum per captured frame.
// If captured frames include multiple places, then frames and detections could
// be distinguished with camera ID.
//
// Usage:
//  ```
//  scouter_multi_region_tracking([tracker_param], [stream],
//                                [instance_states_id_name], [frames_name],
//                                [camera_id_name], [image_name],
//                                [mv_region_name], [timestamp_name])
//  ```
//  [tracker_param]
//  [stream]
//  [instance_states_id_name]
//  [frames_name]
//  [camera_id_name]
//  [image_name]
//  [mv_region_name]
//  [timestamp_name]
//
// Input tuples are required to have following `data.Map` structure, each key
// name is addressed with UDSF's arguments.
//
// Stream Tuple.Data structure:
//  data.Map{
//    "instanceStatesIDFieldName": [ID],
//    "framesFieldName"          : data.Array{
//      []data.Map{
//        "cameraIDFieldName": [camera ID],
//        "imageFiledname"   : [image data] (data.Blob),
//      }
//    },
//    "mvRegionsFieldName": [moving matched detection result] ([]data.Blob),
//    "timestampFieldName": [frame captured time] (data.Timestamp),
//  }
func (c *FramesTrackerStreamFuncCreator) CreateStreamFunction() interface{} {
	return createFramesTrackerUDSF
}

// TypeName returns type name.
func (c *FramesTrackerStreamFuncCreator) TypeName() string {
	return "scouter_multi_region_tracking"
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

package integrator

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type multiPlacesMovingMatcherUDSF struct {
	mvMatcher func(float32,
		[]bridge.RegionsWithCameraID) []bridge.MVCandidate

	integrationIDFieldName data.Path
	aggRegionsFieldName    data.Path
	cameraIDFieldName      data.Path
	regionsFieldName       data.Path
	kthreshold             float32
}

func (sf *multiPlacesMovingMatcherUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {

	integrationID, err := t.Data.Get(sf.integrationIDFieldName)
	if err != nil {
		return err
	}

	aggRegions, err := t.Data.Get(sf.aggRegionsFieldName)
	if err != nil {
		return err
	}
	aggRegionsArray, err := data.AsArray(aggRegions)
	if err != nil {
		return err
	}
	convertedRegions, err := sf.convertToSliceRegions(aggRegionsArray)
	defer func() {
		for _, r := range convertedRegions {
			for _, c := range r.Candidates {
				c.Delete()
			}
		}
	}()
	if err != nil {
		return err
	}

	mvCandidates := sf.mvMatcher(sf.kthreshold, convertedRegions)
	defer func() {
		for _, c := range mvCandidates {
			c.Delete()
		}
	}()

	traceCopyFlag := len(t.Trace) > 0
	for _, c := range mvCandidates {
		now := time.Now()
		m := data.Map{
			"integration_id":         integrationID,
			"moving_matched_regions": data.Blob(c.Serialize()),
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
	return nil
}

// Terminate the components.
func (sf *multiPlacesMovingMatcherUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func (sf *multiPlacesMovingMatcherUDSF) convertToSliceRegions(aggRegions data.Array) (
	[]bridge.RegionsWithCameraID, error) {
	aggRegionsWithID := []bridge.RegionsWithCameraID{}
	for _, regions := range aggRegions {
		regionsMap, err := data.AsMap(regions)
		if err != nil {
			return nil, err
		}
		rWithID, err := sf.lookupRegions(regionsMap)
		if err != nil {
			return nil, err
		}
		aggRegionsWithID = append(aggRegionsWithID, rWithID)
	}
	return aggRegionsWithID, nil
}

func (sf *multiPlacesMovingMatcherUDSF) lookupRegions(regions data.Map) (
	bridge.RegionsWithCameraID, error) {

	empty := bridge.RegionsWithCameraID{}
	id, err := regions.Get(sf.cameraIDFieldName)
	if err != nil {
		return empty, err
	}
	cameraID, err := data.ToInt(id)
	if err != nil {
		return empty, err
	}

	rs, err := regions.Get(sf.regionsFieldName)
	if err != nil {
		return empty, err
	}
	rArray, err := data.AsArray(rs)
	if err != nil {
		return empty, err
	}

	cans := []bridge.Candidate{}
	for _, r := range rArray {
		b, err := data.ToBlob(r)
		if err != nil {
			return empty, err
		}
		candidate := bridge.DeserializeCandidate(b)
		cans = append(cans, candidate)
	}

	return bridge.RegionsWithCameraID{
		CameraID:   int(cameraID),
		Candidates: cans,
	}, nil
}

func createMovingMatcherUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	stream string, integrationIDFieldName string, aggRegionsFieldName string,
	cameraIDFieldName string, regionsFieldName string, kthreshlold float32) (

	udf.UDSF, error) {
	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "moving_matcher",
	}); err != nil {
		return nil, err
	}

	return &multiPlacesMovingMatcherUDSF{
		mvMatcher:              bridge.GetMatching,
		integrationIDFieldName: data.MustCompilePath(integrationIDFieldName),
		aggRegionsFieldName:    data.MustCompilePath(aggRegionsFieldName),
		cameraIDFieldName:      data.MustCompilePath(cameraIDFieldName),
		regionsFieldName:       data.MustCompilePath(regionsFieldName),
		kthreshold:             kthreshlold,
	}, nil
}

// MultiPlacesMovingMatcherUDSFCreator is a creator of multi places moving
// matcher UDSF.
type MultiPlacesMovingMatcherUDSFCreator struct{}

// CreateStreamFunction creates moving matcher stream function for multi places.
//
// Usage:
//  ```
//  scouter_multi_place_moving_matcher([stream], [integration_id_name],
//                                     [agg_regions_name], [camera_id_name],
//                                     [regions_name], [kthreshold])
//  ```
//  [stream]
//  [integration_id_name]
//  [agg_regions_name]
//  [camera_id_name]
//  [regions_name]
//  [kthreshold]
//
// Input tuples are required to have following `data.Map` structure, each key
// name is addressed with UDSF's arguments.
//
// Stream Tuple.Data structure:
//  data.Map{
//    "integrationIDFieldName": [ID],
//    "aggRegionsFildName"    : data.Array{
//      []data.Map{
//        "cameraIDFieldName": [camera ID],
//        "regionsFieldName" : [regions] ([]data.Blob),
//      }
//    }
//  }
func (c *MultiPlacesMovingMatcherUDSFCreator) CreateStreamFunction() interface{} {
	return createMovingMatcherUDSF
}

// TypeName returns type name.
func (c *MultiPlacesMovingMatcherUDSFCreator) TypeName() string {
	return "scouter_multi_place_moving_matcher"
}

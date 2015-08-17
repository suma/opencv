package integrator

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func movingMatcherBatch(ctx *core.Context, multiPlaceRegions data.Array,
	kThreashold float32) (data.Array, error) {

	convertedRegions, err := convertToSliceRegions(multiPlaceRegions)
	defer func() {
		for _, r := range convertedRegions {
			for _, c := range r.Candidates {
				c.Delete()
			}
		}
	}()
	if err != nil {
		return nil, err
	}

	mvCandidates := bridge.GetMatching(kThreashold, convertedRegions)
	defer func() {
		for _, c := range mvCandidates {
			c.Delete()
		}
	}()

	cans := data.Array{}
	for _, c := range mvCandidates {
		b := data.Blob(c.Serialize())
		cans = append(cans, b)
	}
	return cans, nil
}

func convertToSliceRegions(aggRegions data.Array) (
	[]bridge.RegionsWithCameraID, error) {

	aggRegionsWithID := []bridge.RegionsWithCameraID{}
	for _, regions := range aggRegions {
		regionsMap, err := data.AsMap(regions)
		if err != nil {
			return nil, err
		}
		rWithID, err := lookupRegions(regionsMap)
		if err != nil {
			return nil, err
		}
		aggRegionsWithID = append(aggRegionsWithID, rWithID)
	}
	return aggRegionsWithID, nil
}

var regionsPath = data.MustCompilePath("regions")

func lookupRegions(regions data.Map) (bridge.RegionsWithCameraID, error) {

	empty := bridge.RegionsWithCameraID{}
	id, err := regions.Get(cameraIDPath)
	if err != nil {
		return empty, err
	}
	cameraID, err := data.AsInt(id)
	if err != nil {
		return empty, err
	}

	rs, err := regions.Get(regionsPath)
	if err != nil {
		return empty, err
	}
	rArray, err := data.AsArray(rs)
	if err != nil {
		return empty, err
	}

	cans := []bridge.Candidate{}
	for _, r := range rArray {
		b, err := data.AsBlob(r)
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

// MultiPlacesMovingMatcherBatchUDFCreator is a creator of multi places moving
// matcher UDSF.
type MultiPlacesMovingMatcherBatchUDFCreator struct{}

// CreateFunction creates moving matcher batch function for multi places.
// The function will return moving detection result array, the type is
// `[]data.Blob`.
//
// Usage:
//  `multi_place_moving_matcher_batch(multi_place_regions, kThreashold)`
//  multi_place_regions is required following `data.Array`
//    data.Array{
//      []data.Map{
//        "camera_id": [camera ID],
//        "regions"  : [regions] (data.Array of data.Blob),
//      }
//    }
//  kThreshold is the threshold parameter of matching.
func (c *MultiPlacesMovingMatcherBatchUDFCreator) CreateFunction() interface{} {
	return movingMatcherBatch
}

func (c *MultiPlacesMovingMatcherBatchUDFCreator) TypeName() string {
	return "scouter_multi_place_moving_matcher_batch"
}

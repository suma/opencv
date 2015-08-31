package integrator

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func movingMatcherBatch(ctx *core.Context, multiPlaceRegions data.Array,
	kthreashold float32) (data.Array, error) {

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

	mvCandidates := bridge.GetMatching(kthreashold, convertedRegions)
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

func lookupRegions(regions data.Map) (bridge.RegionsWithCameraID, error) {

	empty := bridge.RegionsWithCameraID{}
	id, err := regions.Get(utils.CameraIDPath)
	if err != nil {
		return empty, err
	}
	cameraID, err := data.ToInt(id)
	if err != nil {
		return empty, err
	}

	rs, err := regions.Get(utils.RegionsPath)
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

// MultiPlacesMovingMatcherBatchUDFCreator is a creator of multi places moving
// matcher UDSF.
type MultiPlacesMovingMatcherBatchUDFCreator struct{}

// CreateFunction creates moving matcher batch function for multi places.
// The function will return moving detection result array, the type is
// `[]data.Blob`.
//
// Usage:
//  `scouter_multi_place_moving_matcher_batch([multi_place_regions],
//                                            [kthreshold])`
//  [multi_place_regions]
//    * type: data.Array
//    * regions with camera ID data, required following data.Map structure.
//      data.Array{
//        []data.Map{
//          "camera_id": [camera ID],
//          "regions"  : [regions] (data.Array of data.Blob),
//        }
//      }
//  [kthreshold]
//    * type: float
//    * threshold parameter of matching.
func (c *MultiPlacesMovingMatcherBatchUDFCreator) CreateFunction() interface{} {
	return movingMatcherBatch
}

// TypeName returns type name.
func (c *MultiPlacesMovingMatcherBatchUDFCreator) TypeName() string {
	return "scouter_multi_place_moving_matcher_batch"
}

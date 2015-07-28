package detector

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type ACFDetectBatchFuncCreator struct{}

// CreateFunction returns ACF Detection function. The function will return
// detection result array, the type is `[]data.Blob`.
//
// Usage:
//  `acf_detector_batch('detect_param', frame)`
//    'detect_param' is a parameter name of "acf_detection_parameter" state
//    frame is a captured frame map (`data.Map`), the function required
//      data.Map{
//          "projected_img": [image binary] (`data.Blob`)
//          "offset_x"     : [frame offset x] (`data.Int`)
//          "offset_y"     : [frame offset y] (`data.Int`)
//      }
func (c *ACFDetectBatchFuncCreator) CreateFunction() interface{} {
	return acfDetectBatch
}

func (c *ACFDetectBatchFuncCreator) TypeName() string {
	return "acf_detector_batch"
}

func acfDetectBatch(ctx *core.Context, detectParam string, frame data.Map) (
	data.Array, error) {

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	img, err := lookupFrameData(frame)
	if err != nil {
		return nil, err
	}
	offsetX, offsetY, err := lookupOffsets(frame)
	if err != nil {
		return nil, err
	}

	imgPtr := bridge.DeserializeMatVec3b(img)
	defer imgPtr.Delete()
	candidates := s.d.ACFDetect(imgPtr, offsetX, offsetY)
	defer func() {
		for _, c := range candidates {
			c.Delete()
		}
	}()

	cans := data.Array{}
	for _, c := range candidates {
		b := data.Blob(c.Serialize())
		cans = append(cans, b)
	}
	return cans, nil
}

type FilterByMaskBatchFuncCreator struct{}

func filterByMaskBatch(ctx *core.Context, detectParam string, regions data.Array) (
	data.Array, error) {

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	filteredCans := data.Array{}
	for _, r := range regions {
		regionByte, err := data.AsBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		filter := func() {
			defer regionPtr.Delete()
			if !s.d.FilterByMask(regionPtr) {
				filteredCans = append(filteredCans, r)
			}
		}
		filter()
	}

	return filteredCans, nil
}

// CreateFunction returns filtered by mask function for ACF detection.
// The function will return detection result array, the type is `[]data.Blob`.
//
// Usage:
//  `filter_by_mask_batch('detect_param', regions)`
//    'detect_param' is a parameter name of "acf_detection_parameter" state
//    regions are detection results, which are detected by ACF detector,
//      required `[]data.Blob` type.
func (c *FilterByMaskBatchFuncCreator) CreateFunction() interface{} {
	return filterByMaskBatch
}

func (c *FilterByMaskBatchFuncCreator) TypeName() string {
	return "filter_by_mask_batch"
}

type EstimateHeightBatchFuncCreator struct{}

func estimateHeightBatch(ctx *core.Context, detectParam string, frame data.Map,
	regions data.Array) (data.Array, error) {

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	offsetX, offsetY, err := lookupOffsets(frame)
	if err != nil {
		return nil, err
	}

	estimatedCans := data.Array{}
	for _, r := range regions {
		regionByte, err := data.AsBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		estimate := func() {
			defer regionPtr.Delete()
			s.d.EstimateHeight(&regionPtr, offsetX, offsetY)
			estimatedCans = append(estimatedCans, data.Blob(regionPtr.Serialize()))
		}
		estimate()
	}

	return estimatedCans, nil
}

// CreateFunction returns estimated height function for ACF detection.
// The function will return detection result array, the type is `[]data.Blob`.
//
// Usage:
//  `estimate_height_batch('detect_param', frame, regions)`
//    'detect_param' is a parameter name of "acf_detection_parameter" state
//    frame is a captured frame map (`data.Map`), the function required
//      data.Map{
//          "offset_x"  : [frame offset x] (`data.Int`)
//          "offset_y"  : [frame offset y] (`data.Int`)
//      }
//    regions are detection results, which are detected by ACF detector,
//      required `[]data.Blob` type.
func (c *EstimateHeightBatchFuncCreator) CreateFunction() interface{} {
	return estimateHeightBatch
}

func (c *EstimateHeightBatchFuncCreator) TypeName() string {
	return "estimate_height_batch"
}

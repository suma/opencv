package detector

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type MMDetectBatchFuncCreator struct{}

// CreateFunction returns Multi Model Detection function. The function will
// return detection result array, the type is `[]data.Blob`.
//
// Usage:
//  `mm_detector_batch('detect_param', frame)`
//    'detect_param' is a parameter name of "multi_model_detection_parameter" state
//    frame is a captured frame map (`data.Map`), the function required
//      data.Map{
//          "projected_img": [image binary] (`data.Blob`)
//          "offset_x"     : [frame offset x] (`data.Int`)
//          "offset_y"     : [frame offset y] (`data.Int`)
//      }
func (c *MMDetectBatchFuncCreator) CreateFunction() interface{} {
	return mmDetectBatch
}

func (c *MMDetectBatchFuncCreator) TypeName() string {
	return "mm_detector_batch"
}

func mmDetectBatch(ctx *core.Context, detectParam string, frame data.Map) (
	data.Array, error) {

	s, err := lookupMMDetectParamState(ctx, detectParam)
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
	candidates := s.d.MMDetect(imgPtr, offsetX, offsetY)
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

type FilterByMaskMMBatchFuncCreator struct{}

func filterByMaskMMBatch(ctx *core.Context, detectParam string, regions data.Array) (
	data.Array, error) {

	s, err := lookupMMDetectParamState(ctx, detectParam)
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

// CreateFunction returns filtered by mask function for multi-model detection.
// The function will return detection result array, the type is `[]data.Blob`.
//
// Usage:
//  `multi_model_filter_by_mask_batch('detect_param', regions)`
//    'detect_param' is a parameter name of "multi_model_detection_parameter" state
//    regions are detection results, which are detected by multi-model detector,
//      required `[]data.Blob` type.
func (c *FilterByMaskMMBatchFuncCreator) CreateFunction() interface{} {
	return filterByMaskBatch
}

func (c *FilterByMaskMMBatchFuncCreator) TypeName() string {
	return "multi_model_filter_by_mask_batch"
}

type EstimateHeightMMBatchFuncCreator struct{}

func estimateHeightMMBatch(ctx *core.Context, detectParam string, frame data.Map,
	regions data.Array) (data.Array, error) {

	s, err := lookupMMDetectParamState(ctx, detectParam)
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

// CreateFunction returns estimated height function for multi-model detection.
// The function will return detection result array, the type is `[]data.Blob`.
//
// Usage:
//  `multi_model_estimate_height_batch('detect_param', frame, regions)`
//    'detect_param' is a parameter name of "multi_model_detection_parameter" state
//    frame is a captured frame map (`data.Map`), the function required
//      data.Map{
//          "offset_x"  : [frame offset x] (`data.Int`)
//          "offset_y"  : [frame offset y] (`data.Int`)
//      }
//    regions are detection results, which are detected by multi-model detector,
//      required `[]data.Blob` type.
func (c *EstimateHeightMMBatchFuncCreator) CreateFunction() interface{} {
	return estimateHeightBatch
}

func (c *EstimateHeightMMBatchFuncCreator) TypeName() string {
	return "multi_model_estimate_height_batch"
}

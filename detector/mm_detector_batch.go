package detector

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// MMDetectBatchFuncCreator is a creator of multi-model detector UDF.
type MMDetectBatchFuncCreator struct{}

// CreateFunction returns Multi Model Detection function.
//
// Usage:
//  `scouter_mm_detector_batch([detect_param], [frame])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_mm_detection_param" state
//  [frame]
//    * type: data.Map
//    * captured frame which are applied `scouter_frame_applier` UDF. The
//      frame's map structure is required following structure.
//      data.Map{
//        "projected_img": [image binary] (`data.Blob`)
//        "offset_x":      [frame offset x] (`data.Int`)
//        "offset_y":      [frame offset y] (`data.Int`)
//      }
//
// Return:
//  The function will return detected regions array, the type is `[]data.Blob`.
func (c *MMDetectBatchFuncCreator) CreateFunction() interface{} {
	return mmDetectBatch
}

// TypeName returns type name.
func (c *MMDetectBatchFuncCreator) TypeName() string {
	return "scouter_mm_detector_batch"
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

	cans := make(data.Array, len(candidates))
	for i, c := range candidates {
		b := data.Blob(c.Serialize())
		cans[i] = b
	}
	return cans, nil
}

// FilterByMaskMMBatchFuncCreator is a creator of filtering by bask UDF.
type FilterByMaskMMBatchFuncCreator struct{}

func filterByMaskMMBatch(ctx *core.Context, detectParam string, regions data.Array) (
	data.Array, error) {

	s, err := lookupMMDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	// filterCans size is not same as len(regions), and not use make()
	filteredCans := data.Array{}
	for _, r := range regions {
		regionByte, err := data.ToBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		func() {
			defer regionPtr.Delete()
			if !s.d.FilterByMask(regionPtr) {
				filteredCans = append(filteredCans, r)
			}
		}()
	}

	return filteredCans, nil
}

// CreateFunction creates a batch filter by mask for multi model detection.
//
// Usage:
//  `scouter_mm_filter_by_mask_batch([detect_param], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_mm_detection_param" state
//  [regions]
//    * type: []data.Blob
//    * detected regions, which are applied multi model detection.
//
// Returns:
//  The function will return filtered regions array, the type is `[]data.Blob`.
func (c *FilterByMaskMMBatchFuncCreator) CreateFunction() interface{} {
	return filterByMaskBatch
}

// TypeName returns type name.
func (c *FilterByMaskMMBatchFuncCreator) TypeName() string {
	return "scouter_mm_filter_by_mask_batch"
}

// EstimateHeightMMBatchFuncCreator is creator of height estimator UDF.
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

	estimatedCans := make(data.Array, len(regions))
	for i, r := range regions {
		regionByte, err := data.ToBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		func() {
			defer regionPtr.Delete()
			s.d.EstimateHeight(&regionPtr, offsetX, offsetY)
			estimatedCans[i] = data.Blob(regionPtr.Serialize())
		}()
	}

	return estimatedCans, nil
}

// CreateFunction creates a estimate height function for multi model detection.
//
// Usage:
//  `scouter_mm_estimate_height_batch([detect_param], [frame], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_mm_detection_param" state
//  [frame]
//    * type: data.Map
//    * captured frame which are applied `scouter_frame_applier` UDF. The
//      frame's map structure is required following structure.
//      data.Map{
//        "offset_x"  : [frame offset x] (`data.Int`)
//        "offset_y"  : [frame offset y] (`data.Int`)
//      }
//  [regions]
//    * type: []data.Blob
//    * detected regions, which are applied multi model detection.
//    * these regions are detected from [frame]
//
// Return:
//   The function will return estimate regions array, the type is `[]data.Blob`.
func (c *EstimateHeightMMBatchFuncCreator) CreateFunction() interface{} {
	return estimateHeightBatch
}

// TypeName returns type name.
func (c *EstimateHeightMMBatchFuncCreator) TypeName() string {
	return "scouter_mm_estimate_height_batch"
}

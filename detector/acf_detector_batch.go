package detector

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// ACFDetectBatchFuncCreator is a creator of ACF detector UDF.
type ACFDetectBatchFuncCreator struct{}

// CreateFunction returns ACF Detection function.
//
// Usage:
//  `scouter_acf_detector_batch([detect_param], [frame], [camera_id])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_acf_detection_param" state
//  [frame]
//    * type: data.Map
//    * captured frame which are applied `scouter_frame_applier` UDF. The
//      frame's map structure is required following structure.
//      data.Map{
//        "projected_img": [image binary] (`data.Blob`)
//        "offset_x":      [frame offset x] (`data.Int`)
//        "offset_y":      [frame offset y] (`data.Int`)
//      }
//  [camera ID]
//    * type: int
//    * camera ID, to use for detection result, not use for detection.
//
// Return:
//  The function will return detected regions array, the type is `[]data.Blob`.
func (c *ACFDetectBatchFuncCreator) CreateFunction() interface{} {
	return acfDetectBatch
}

// TypeName returns type name.
func (c *ACFDetectBatchFuncCreator) TypeName() string {
	return "scouter_acf_detector_batch"
}

func acfDetectBatch(ctx *core.Context, detectParam string, frame data.Map,
	cameraID int) (data.Array, error) {

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
	candidates := s.d.ACFDetect(imgPtr, offsetX, offsetY, cameraID)
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

// FilterByMaskBatchFuncCreator is a creator of filtering by bask UDF.
type FilterByMaskBatchFuncCreator struct{}

func filterByMaskBatch(ctx *core.Context, detectParam string, regions data.Array) (
	data.Array, error) {

	s, err := lookupACFDetectParamState(ctx, detectParam)
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

// CreateFunction creates a batch filter by mask for ACF detection.
//
// Usage:
//  `scouter_filter_by_mask_batch([detect_param], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_acf_detection_param" state
//  [regions]
//    * type: []data.Blob
//    * detected regions, which are applied ACF detection.
//
// Returns:
//  The function will return filtered regions array, the type is `[]data.Blob`.
func (c *FilterByMaskBatchFuncCreator) CreateFunction() interface{} {
	return filterByMaskBatch
}

// TypeName returns type name.
func (c *FilterByMaskBatchFuncCreator) TypeName() string {
	return "scouter_filter_by_mask_batch"
}

// EstimateHeightBatchFuncCreator is creator of height estimator UDF.
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

// CreateFunction creates a estimate height function for ACF detection.
//
// Usage:
//  `estimate_height_batch([detect_param], [frame], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_acf_detection_param" state
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
//    * detected regions, which are applied ACF detection.
//    * these regions are detected from [frame]
//
// Return:
//   The function will return estimate regions array, the type is `[]data.Blob`.
func (c *EstimateHeightBatchFuncCreator) CreateFunction() interface{} {
	return estimateHeightBatch
}

// TypeName returns type name.
func (c *EstimateHeightBatchFuncCreator) TypeName() string {
	return "scouter_estimate_height_batch"
}

// PutFeatureBatchUDFCreator is a creator of putting feature.
type PutFeatureBatchUDFCreator struct{}

func putFeatureBatch(ctx *core.Context, detectParam string, image []byte,
	regions data.Array) (data.Array, error) {

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	imgPtr := bridge.DeserializeMatVec3b(image)
	defer imgPtr.Delete()

	putFeatureCans := make(data.Array, len(regions))
	for i, r := range regions {
		regionByte, err := data.ToBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		func() {
			defer regionPtr.Delete()
			s.d.PutFeature(&regionPtr, imgPtr)
			putFeatureCans[i] = data.Blob(regionPtr.Serialize())
		}()
	}

	return putFeatureCans, nil
}

// CreateFunction create a putting feature function for ACF detection.
//
// Usage:
//  `scouter_put_feature_batch([detect_param], [image], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_acf_detection_param" state
//  [image]
//    * type: []byte
//    * captured image
//  [regions]
//    * type: []data.Blob
//    * detected regions, which are applied ACF detection.
//    * these regions are detected from [frame]
//
// return:
//  The function will return regions array which regions is set features, the
//  type is `[]data.Blob`
func (c *PutFeatureBatchUDFCreator) CreateFunction() interface{} {
	return putFeatureBatch
}

// TypeName returns type name
func (c *PutFeatureBatchUDFCreator) TypeName() string {
	return "scouter_put_feature_batch"
}

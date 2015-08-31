package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// FrameApplierFuncCreator is a creator of frame processing applier UDF.
type FrameApplierFuncCreator struct{}

func frameApplier(ctx *core.Context, fpParam string, capture []byte) (
	data.Map, error) {
	s, err := lookupFPParamState(ctx, fpParam)
	if err != nil {
		return nil, err
	}

	bufp := bridge.DeserializeMatVec3b(capture)
	defer bufp.Delete()
	img, offsetX, offsetY := s.fp.Projection(bufp)

	m := data.Map{
		"projected_img": data.Blob(img.Serialize()),
		"offset_x":      data.Int(offsetX),
		"offset_y":      data.Int(offsetY),
	}

	return m, nil
}

// CreateFunction crates frame processing applier UDF.
//
// Usage:
//  `scouter_frame_applier([frame_processor_param_state], [captured image])`
//  [frame_processor_param_state]
//    * type: string
//    * the state name of frame processor parameter
//  [captured image]
//    * type: []byte serialized from `cv::Mat_<cv::Vec3b>`
//    * the image data
//
// Return:
//  The function will return following `data.Map`.
//  data.Map{
//    "projected_img": [the result of projected image] (data.Blob)
//                     (serialized from `cv::Mat_<cv::Vec3b>`),
//    "offset_x":      [the value of offset X axis] (data.Int),
//    "offset_y":      [the value of offset Y axis] (data.Int),
//  }
func (c *FrameApplierFuncCreator) CreateFunction() interface{} {
	return frameApplier
}

// TypeName returns type name.
func (c *FrameApplierFuncCreator) TypeName() string {
	return "scouter_frame_applier"
}

func lookupFPParamState(ctx *core.Context, stateName string) (
	*FrameProcessorParamState, error) {
	st, err := ctx.SharedStates.Get(stateName)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*FrameProcessorParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to frame_processor_parameter.state",
		stateName)
}

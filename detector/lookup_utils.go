package detector

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

var (
	projectedIMGPath = data.MustCompilePath("projected_img")
	offsetXPath      = data.MustCompilePath("offset_x")
	offsetYPath      = data.MustCompilePath("offset_y")
)

func lookupACFDetectParamState(ctx *core.Context, detectParam string) (
	*ACFDetectionParamState, error) {
	st, err := ctx.SharedStates.Get(detectParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ACFDetectionParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to acf_detection_param.state",
		detectParam)
}

func lookupMMDetectParamState(ctx *core.Context, detectParam string) (
	*MMDetectionParamState, error) {
	st, err := ctx.SharedStates.Get(detectParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*MMDetectionParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to mm_detection_param.state", detectParam)
}

func lookupFrameData(frame data.Map) ([]byte, error) {
	img, err := frame.Get(projectedIMGPath)
	if err != nil {
		return []byte{}, err
	}
	image, err := data.AsBlob(img)
	if err != nil {
		return []byte{}, err
	}

	return image, nil
}

func lookupOffsets(frame data.Map) (int, int, error) {
	ox, err := frame.Get(offsetXPath)
	if err != nil {
		return 0, 0, err
	}
	offsetX, err := data.AsInt(ox)
	if err != nil {
		return 0, 0, err
	}

	oy, err := frame.Get(offsetYPath)
	if err != nil {
		return 0, 0, err
	}
	offsetY, err := data.AsInt(oy)
	if err != nil {
		return 0, 0, err
	}

	return int(offsetX), int(offsetY), nil
}

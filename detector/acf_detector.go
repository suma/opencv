package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

func ACFDetectFunc(ctx *core.Context, detectParam string, frame tuple.Map) (tuple.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	img, err := lookupFrameData(frame)
	if err != nil {
		return nil, err
	}
	offsetX, offsetY, err := loopupOffsets(frame)
	if err != nil {
		return nil, err
	}
	candidates := s.d.ACFDetect(bridge.DeserializeMatVec3b(img), offsetX, offsetY)
	detected := tuple.Array{}
	for _, candidate := range candidates {
		detected = append(detected, tuple.Blob(candidate.Serialize()))
		candidate.Delete() // TODO use defer
	}
	frame["detect"] = detected
	return frame, nil
}

func FilterByMaskFunc(ctx *core.Context, detectParam string, frame tuple.Map) (tuple.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	candidates, err := frame.Get("detect")
	if err != nil {
		return nil, err
	}
	cans, err := tuple.AsArray(candidates)
	if err != nil {
		return nil, err
	}

	cansByte := [][]byte{}
	for _, c := range cans {
		b, err := tuple.AsBlob(c)
		if err != nil {
			return nil, err // TODO return is OK?
		}
		cansByte = append(cansByte, b)
	}

	filteredCans := s.d.FilterByMask(cansByte)
	filtered := tuple.Array{}
	for _, fc := range filteredCans {
		filtered = append(filtered, tuple.Blob(fc.Serialize()))
		fc.Delete() // TODO use defer
	}

	frame["detect"] = filtered // TODO overwrite is OK?
	return frame, nil
}

func EstimateHeight(ctx *core.Context, detectParam string, frame tuple.Map) (tuple.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	offsetX, offsetY, err := loopupOffsets(frame)
	if err != nil {
		return nil, err
	}

	candidates, err := frame.Get("detect")
	if err != nil {
		return nil, err
	}
	cans, err := tuple.AsArray(candidates)
	if err != nil {
		return nil, err
	}

	cansByte := [][]byte{}
	for _, c := range cans {
		b, err := tuple.AsBlob(c)
		if err != nil {
			return nil, err // TODO return is OK?
		}
		cansByte = append(cansByte, b)
	}

	estimatedCans := s.d.EstimateHeight(cansByte, offsetX, offsetY)
	estimated := tuple.Array{}
	for _, ec := range estimatedCans {
		estimated = append(estimated, tuple.Blob(ec.Serialize()))
		ec.Delete() // TODO use defer
	}

	frame["detect"] = estimated // TODO overwrite is OK?
	return estimated, nil
}

func lookupACFDetectParamState(ctx *core.Context, detectParam string) (*ACFDetectionParamState, error) {
	st, err := ctx.GetSharedState(detectParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ACFDetectionParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to acf_detection_param.state", detectParam)
}

func lookupFrameData(frame tuple.Map) ([]byte, error) {
	img, err := frame.Get("projected_img")
	if err != nil {
		return []byte{}, err
	}
	image, err := tuple.AsBlob(img)
	if err != nil {
		return []byte{}, err
	}

	return image, nil
}

func loopupOffsets(frame tuple.Map) (int, int, error) {
	ox, err := frame.Get("offset_x")
	if err != nil {
		return 0, 0, err
	}
	offsetX, err := tuple.AsInt(ox)
	if err != nil {
		return 0, 0, err
	}

	oy, err := frame.Get("offset_y")
	if err != nil {
		return 0, 0, err
	}
	offsetY, err := tuple.AsInt(oy)
	if err != nil {
		return 0, 0, err
	}

	return int(offsetX), int(offsetY), nil
}

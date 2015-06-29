package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

func ACFDetectFunc(ctx *core.Context, detectParam tuple.Value, frame tuple.Value) (tuple.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	img, offsetX, offsetY, err := lookupFrameData(ctx, frame)
	if err != nil {
		return nil, err
	}
	candidates := s.d.ACFDetect(bridge.DeserializeMatVec3b(img), offsetX, offsetY)
	detected := tuple.Array{}
	for _, candidate := range candidates {
		detected = append(detected, tuple.Blob(candidate.Serialize()))
	}
	return detected, nil
}

func lookupACFDetectParamState(ctx *core.Context, stateName tuple.Value) (*ACFDetectionParamState, error) {
	name, err := tuple.AsString(stateName)
	if err != nil {
		return nil, fmt.Errorf("name of the state must be a string: %v", stateName)
	}

	st, err := ctx.GetSharedState(name)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ACFDetectionParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to acf_detection_param.state", name)
}

func lookupFrameData(ctx *core.Context, frameValue tuple.Value) ([]byte, int, int, error) {
	frame, err := tuple.AsMap(frameValue)
	if err != nil {
		return []byte{}, 0, 0, err
	}

	img, err := frame.Get("projected_img")
	if err != nil {
		return []byte{}, 0, 0, err
	}
	image, err := tuple.AsBlob(img)
	if err != nil {
		return []byte{}, 0, 0, err
	}

	ox, err := frame.Get("offset_x")
	if err != nil {
		return []byte{}, 0, 0, err
	}
	offset_x, err := tuple.AsInt(ox)
	if err != nil {
		return []byte{}, 0, 0, err
	}

	oy, err := frame.Get("offset_y")
	if err != nil {
		return []byte{}, 0, 0, err
	}
	offset_y, err := tuple.AsInt(oy)
	if err != nil {
		return []byte{}, 0, 0, err
	}
	return image, int(offset_x), int(offset_y), nil
}

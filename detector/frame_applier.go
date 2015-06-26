package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

func Func(ctx *core.Context, cameraParam tuple.Value, captureMat tuple.Value) (tuple.Value, error) {
	s, err := lookupCamerparameState(ctx, cameraParam)
	if err != nil {
		return nil, err
	}

	capMat, err := tuple.AsMap(captureMat)
	if err != nil {
		return nil, fmt.Errorf("capture data must be a Map: %v", err.Error())
	}
	buf, _, err := lookupMatData(ctx, capMat)
	if err != nil {
		return nil, err
	}

	img, offsetX, offsetY := s.fp.Projection(bridge.DeserializeMatVec3b(buf))
	capMat["projection_img"] = tuple.Blob(img.Serialize())
	capMat["offset_x"] = tuple.Int(offsetX)
	capMat["offset_y"] = tuple.Int(offsetY)
	return capMat, nil
}

func lookupCamerparameState(ctx *core.Context, stateName tuple.Value) (*CameraParameterState, error) {
	name, err := tuple.AsString(stateName)
	if err != nil {
		return nil, fmt.Errorf("name of the state must be a string: %v", stateName)
	}

	st, err := ctx.GetSharedState(name)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*CameraParameterState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to camera_parameter.state", name)
}

func lookupMatData(ctx *core.Context, capMat tuple.Map) ([]byte, int64, error) {
	mat, err := capMat.Get("capture")
	if err != nil {
		return []byte{}, 0, err
	}
	buf, err := tuple.AsBlob(mat)
	if err != nil {
		return []byte{}, 0, err
	}

	ci, err := capMat.Get("cameraID")
	if err != nil {
		return []byte{}, 0, err
	}
	cameraID, err := tuple.AsInt(ci)
	if err != nil {
		return []byte{}, 0, err
	}
	return buf, cameraID, nil
}

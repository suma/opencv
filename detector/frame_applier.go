package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func FrameApplierFunc(ctx *core.Context, cameraParam data.Value, captureMat data.Value) (data.Value, error) {
	s, err := lookupCameraParamState(ctx, cameraParam)
	if err != nil {
		return nil, err
	}

	capMat, err := data.AsMap(captureMat)
	if err != nil {
		return nil, fmt.Errorf("capture data must be a Map: %v", err.Error())
	}
	buf, _, err := lookupMatData(ctx, capMat)
	if err != nil {
		return nil, err
	}

	img, offsetX, offsetY := s.fp.Projection(bridge.DeserializeMatVec3b(buf))
	capMat["projected_img"] = data.Blob(img.Serialize())
	capMat["offset_x"] = data.Int(offsetX)
	capMat["offset_y"] = data.Int(offsetY)
	return capMat, nil
}

func lookupCameraParamState(ctx *core.Context, stateName data.Value) (*CameraParamState, error) {
	name, err := data.AsString(stateName)
	if err != nil {
		return nil, fmt.Errorf("name of the state must be a string: %v", stateName)
	}

	st, err := ctx.GetSharedState(name)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*CameraParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to camera_parameter.state", name)
}

func lookupMatData(ctx *core.Context, capMat data.Map) ([]byte, int64, error) {
	mat, err := capMat.Get("capture")
	if err != nil {
		return []byte{}, 0, err
	}
	buf, err := data.AsBlob(mat)
	if err != nil {
		return []byte{}, 0, err
	}

	ci, err := capMat.Get("cameraID")
	if err != nil {
		return []byte{}, 0, err
	}
	cameraID, err := data.AsInt(ci)
	if err != nil {
		return []byte{}, 0, err
	}
	return buf, cameraID, nil
}

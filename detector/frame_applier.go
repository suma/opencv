package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func FrameApplierFunc(ctx *core.Context, cameraParam string, capture data.Blob) (data.Value, error) {
	s, err := lookupCameraParamState(ctx, cameraParam)
	if err != nil {
		return nil, err
	}

	buf, err := data.AsBlob(capture)
	if err != nil {
		return nil, err
	}

	bufp := bridge.DeserializeMatVec3b(buf)
	defer bufp.Delete()
	img, offsetX, offsetY := s.fp.Projection(bufp)

	m := data.Map{
		"projected_img": data.Blob(img.Serialize()),
		"offset_x":      data.Int(offsetX),
		"offset_y":      data.Int(offsetY),
	}

	return m, nil
}

func lookupCameraParamState(ctx *core.Context, stateName string) (*CameraParamState, error) {
	st, err := ctx.SharedStates.Get(stateName)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*CameraParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to camera_parameter.state", stateName)
}

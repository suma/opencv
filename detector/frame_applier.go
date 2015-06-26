package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"time"
)

func Func(ctx *core.Context, cameraParam tuple.Value, captureMat tuple.Value) (tuple.Value, error) {
	s, err := lookupCamerparameState(ctx, cameraParam)
	if err != nil {
		return nil, err
	}

	buf, cameraID, err := lookupMatData(ctx, captureMat)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	inow := now.UnixNano() / int64(time.Millisecond)                 // [ms]
	s.fp.Apply(bridge.DeserializeMatVec3b(buf), inow, int(cameraID)) // TODO now is wrong, should be get capture time
	return nil, nil
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

func lookupMatData(ctx *core.Context, captureMat tuple.Value) ([]byte, int64, error) {
	capMat, err := tuple.AsMap(captureMat)
	if err != nil {
		return []byte{}, 0, fmt.Errorf("capture data must be a Map: %v", err.Error())
	}

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

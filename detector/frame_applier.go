package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type FrameApplierFuncCreator struct{}

func frameApplier(ctx *core.Context, fpParam string, capture data.Blob) (
	data.Value, error) {
	s, err := lookupFPParamState(ctx, fpParam)
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

func (c *FrameApplierFuncCreator) CreateFunction() interface{} {
	return frameApplier
}

func (c *FrameApplierFuncCreator) TypeName() string {
	return "frame_applier"
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

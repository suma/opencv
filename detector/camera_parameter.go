package detector

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

// CameraParameter is used for framing projection.
type CameraParameter struct {
}

func (c *CameraParameter) NewState(ctx *core.Context, with tuple.Map) (core.SharedState, error) {
	return nil, nil
}

func (c *CameraParameter) TypeName() string {
	return "camera_parameter"
}

func (c *CameraParameter) Init(ctx *core.Context) error {
	return nil
}

func (c *CameraParameter) Write(ctx *core.Context, t *tuple.Tuple) error {
	return nil
}

func (c *CameraParameter) Terminate(ctx *core.Context) error {
	return nil
}

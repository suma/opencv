package recog

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
)

func lookupImageTaggerCaffeParamState(ctx *core.Context, taggerParam string) (
	*ImageTaggerCaffeParamState, error) {
	st, err := ctx.SharedStates.Get(taggerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ImageTaggerCaffeParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to image_tagger_caffe_param.state",
		taggerParam)
}

package recog

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func RecognizeFunc(ctx *core.Context, taggerParam string, frame data.Blob, regions data.Array) (data.Value, error) {
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	b, err := data.AsBlob(frame)
	if err != nil {
		return nil, err
	}
	img := bridge.DeserializeMatVec3b(b)

	candidates := []bridge.Candidate{}
	for _, r := range regions {
		b, err := data.AsBlob(r)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, bridge.DeserializeCandidate(b))
	}

	recognized := s.tagger.PredictTagsBatch(candidates, img)
	ret := data.Array{}
	for _, r := range recognized {
		ret = append(ret, data.Blob(r.Serialize()))
		r.Delete() // TODO use defer
	}

	return data.Array(ret), nil
}

func lookupImageTaggerCaffeParamState(ctx *core.Context, taggerParam string) (*ImageTaggerCaffeParamState, error) {
	st, err := ctx.GetSharedState(taggerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ImageTaggerCaffeParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to image_tagger_caffe_param.state", taggerParam)
}

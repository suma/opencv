package recog

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func CropFunc(ctx *core.Context, taggerParam string, region data.Blob, image data.Blob) (data.Value, error) {
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	regionByte, err := data.AsBlob(region)
	if err != nil {
		return nil, err
	}
	r := bridge.DeserializeCandidate(regionByte)

	imageByte, err := data.AsBlob(image)
	if err != nil {
		return nil, err
	}
	img := bridge.DeserializeMatVec3b(imageByte)

	cropped := s.tagger.Crop(r, img)
	return data.Blob(cropped.Serialize()), nil
}

func PredictTagsBatchFunc(ctx *core.Context, taggerParam string, regions data.Array, croppedImgs data.Array) (data.Value, error) {
	if len(regions) != len(croppedImgs) {
		return nil, fmt.Errorf("region size and cropped image size must same [region: %d, cropped image: %d",
			len(regions), len(croppedImgs))
	}
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	candidates := []bridge.Candidate{}
	for _, r := range regions {
		b, err := data.AsBlob(r)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, bridge.DeserializeCandidate(b))
	}

	cropps := []bridge.MatVec3b{}
	for _, c := range croppedImgs {
		b, err := data.AsBlob(c)
		if err != nil {
			return nil, err
		}
		cropps = append(cropps, bridge.DeserializeMatVec3b(b))
	}

	recognized := s.tagger.PredictTagsBatch(candidates, cropps)
	ret := data.Array{}
	for _, r := range recognized {
		ret = append(ret, data.Blob(r.Serialize()))
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

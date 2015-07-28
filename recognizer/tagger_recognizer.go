package recog

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type PredictTagsFuncCreator struct{}

func predictTags(ctx *core.Context, taggerParam string, cropImg []byte, region []byte) ([]byte, error) {
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	croppedImg := bridge.DeserializeMatVec3b(cropImg)
	defer croppedImg.Delete()
	candidate := bridge.DeserializeCandidate(region)
	defer candidate.Delete()

	taggedRegion := s.tagger.PredictTags(candidate, croppedImg)
	defer taggedRegion.Delete()

	return taggedRegion.Serialize(), nil
}

func (c *PredictTagsFuncCreator) CreateFunction() interface{} {
	return predictTags
}

func (c *PredictTagsFuncCreator) TypeName() string {
	return "predict_tags"
}

type CroppingAndPredictTagsFuncCreator struct{}

func croppingAndPredictTags(ctx *core.Context, taggerParam string, region []byte,
	img []byte) ([]byte, error) {
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	image := bridge.DeserializeMatVec3b(img)
	defer image.Delete()
	candidate := bridge.DeserializeCandidate(region)
	defer candidate.Delete()

	taggedRegion := s.tagger.CroppingAndPredictTags(candidate, image)
	defer taggedRegion.Delete()

	return taggedRegion.Serialize(), nil
}

func (c *CroppingAndPredictTagsFuncCreator) CreateFunction() interface{} {
	return croppingAndPredictTags
}

func (c *CroppingAndPredictTagsFuncCreator) TypeName() string {
	return "cropping_and_predict_tags"
}

type DrawDeteciontResultFuncCreator struct{}

func drawDetectionResult(ctx *core.Context, frame []byte, regions data.Array) ([]byte, error) {
	img := bridge.DeserializeMatVec3b(frame)
	defer img.Delete()

	canObjs := []bridge.Candidate{}
	for _, c := range regions {
		b, err := data.AsBlob(c)
		if err != nil {
			return nil, err
		}
		canObjs = append(canObjs, bridge.DeserializeCandidate(b))
	}
	defer func() {
		for _, c := range canObjs {
			c.Delete()
		}
	}()

	ret := bridge.DrawDetectionResultWithTags(img, canObjs)
	defer ret.Delete()
	return ret.Serialize(), nil
}

func (c *DrawDeteciontResultFuncCreator) CreateFunction() interface{} {
	return drawDetectionResult
}

func (c *DrawDeteciontResultFuncCreator) TypeName() string {
	return "draw_detection_result_with_tags"
}

package recog

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// CroppingAndPredictTagsFuncCreator is a creator of cropping and predict tags
// UDF.
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

	taggedRegion := s.tagger.CropAndPredictTags(candidate, image)
	defer taggedRegion.Delete()

	return taggedRegion.Serialize(), nil
}

// CreateFunction returns cropping and predict tags function. This function
// executes two tasks. First, cropping an image took by tagger parameters.
// Second, predicting tags and return region with the tags. Tags information is
// set in a region.
//
// Usage:
//  `scouter_crop_and_predict_tags([tagger_param], [region], [image])`
//  [tagger_param]
//    * type: string
//    * a parameter name of "scouter_image_tagger_caffe" state
//  [region]
//    * type: []byte
//    * a detected region created by detected function.
//    * these regions are detected from [image]
//  [image]
//    * type: []byte
//    * a captured image
// Return:
//  The function will return a tagging region, the type is `[]byte`
func (c *CroppingAndPredictTagsFuncCreator) CreateFunction() interface{} {
	return croppingAndPredictTags
}

// TypeName returns type name.
func (c *CroppingAndPredictTagsFuncCreator) TypeName() string {
	return "scouter_crop_and_predict_tags"
}

// DrawDeteciontResultFuncCreator is a creator of drawing regions with tags in a
// frame UDF.
type DrawDeteciontResultFuncCreator struct{}

func drawDetectionResult(ctx *core.Context, frame []byte, regions data.Array) (
	[]byte, error) {
	img := bridge.DeserializeMatVec3b(frame)
	defer img.Delete()

	canObjs := make([]bridge.Candidate, len(regions))
	for i, c := range regions {
		b, err := data.AsBlob(c)
		if err != nil {
			return nil, err
		}
		canObjs[i] = bridge.DeserializeCandidate(b)
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

// CreateFunction creates a drawing regions with tags on a frame function.
//
// Usage:
//  `scouter_draw_regions_with_tags([frame], [regions])`
//  [frame]
//    * type: []byte
//    * captured frame, which is serialized from `cv::Mat_<cv::Vec3b>`.
//  [regions]
//    * type: []data.Blob
//    * detected regions, which are applied prediction UDF/UDSF
//    * these regions are detected from [frame]
//
// Return:
//  The function will return an image data serialized from `cv::Mat_<cv::Vec3b>`,
//  the type is `[]byte`
func (c *DrawDeteciontResultFuncCreator) CreateFunction() interface{} {
	return drawDetectionResult
}

// TypeName returns type name.
func (c *DrawDeteciontResultFuncCreator) TypeName() string {
	return "scouter_draw_regions_with_tags"
}

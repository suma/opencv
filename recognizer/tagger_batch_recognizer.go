package recog

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type predictTagsBatchUDSF struct {
	predictTagsBatch func([]bridge.Candidate, bridge.MatVec3b) []bridge.Candidate
	frameIDName      data.Path
	regionsName      data.Path
	imageName        data.Path
}

// Process streams tagged regions, which is serialized from
// `scouter::ObjectCandidate`. Tags information is set in a region.
//
// Stream Tuple.Data structure:
//  data.Map{
//    "frame_id":      [frame ID] (`data.Int`),
//    "regions_count": [size of total tagged regions in a frame] (`data.Int`),
//    "region":        [tagged region] (`data.Blob`),
//  }
func (sf *predictTagsBatchUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {
	frameID, err := t.Data.Get(sf.frameIDName)
	if err != nil {
		return err
	}

	regionsData, err := t.Data.Get(sf.regionsName)
	if err != nil {
		return err
	}
	regions, err := data.AsArray(regionsData)
	if err != nil {
		return err
	}

	imgData, err := t.Data.Get(sf.imageName)
	if err != nil {
		return err
	}
	img, err := data.ToBlob(imgData)
	if err != nil {
		return err
	}

	candidates := make([]bridge.Candidate, len(regions))
	defer func() {
		for _, c := range candidates {
			c.Delete()
		}
	}()
	for i, r := range regions {
		rb, err := data.ToBlob(r)
		if err != nil {
			return err
		}
		candidates[i] = bridge.DeserializeCandidate(rb)
	}

	imgPtr := bridge.DeserializeMatVec3b(img)
	defer imgPtr.Delete()

	recognized := sf.predictTagsBatch(candidates, imgPtr)
	defer func() {
		for _, r := range recognized {
			r.Delete()
		}
	}()

	traceCopyFlag := len(t.Trace) > 0
	for _, r := range recognized {
		now := time.Now()
		m := data.Map{
			"frame_id":           frameID,
			"regions_count":      data.Int(len(recognized)),
			"region_with_tagger": data.Blob(r.Serialize()),
		}
		traces := []core.TraceEvent{}
		if traceCopyFlag { // reduce copy cost when trace mode is off
			traces = make([]core.TraceEvent, len(t.Trace), (cap(t.Trace)+1)*2)
			copy(traces, t.Trace)
		}
		tu := &core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: t.ProcTimestamp,
			Trace:         traces,
		}
		w.Write(ctx, tu)
	}
	return nil
}

// Terminate the components.
func (sf *predictTagsBatchUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createPredictTagsBatchUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	taggerParam string, stream string, frameIDName string, regionsName string,
	imageName string) (udf.UDSF, error) {
	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "scouter_predict_tags_batch_stream",
	}); err != nil {
		return nil, err
	}

	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	return &predictTagsBatchUDSF{
		predictTagsBatch: s.tagger.CropAndPredictTagsBatch,
		frameIDName:      data.MustCompilePath(frameIDName),
		regionsName:      data.MustCompilePath(regionsName),
		imageName:        data.MustCompilePath(imageName),
	}, nil
}

// PredictTagsBatchStreamFuncCreator is a creator of predicting tags UDSF.
type PredictTagsBatchStreamFuncCreator struct{}

// CreateStreamFunction returns predicting tags stream function. This stream
// function requires ID per frame to determine the regions detected from.
//
// Usage:
//  ```
//  scouter_predict_tags_batch_stream([tagger_param], [stream], [frame_id_name],
//                                    [regions_name], [image_name])
//  ```
//  [tagger_param]
//    * type: string
//    * a parameter name of "scouter_image_tagger_caffe" state
//  [stream]
//    * type: string
//    * a input stream name, see following stream spec.
//  [frame_id_name]
//    * type: string
//    * a field name of Frame ID
//    * if empty then applied "frame_id"
//  [regions_name]
//    * type: string
//    * a field name of regions
//    * if empty, then applied "regions"
//  [image_name]
//    * type: string
//    * a file name of captured image
//    * if empty, then applied "image"
//
// Input tuples are required to have following `data.Map` structures. The keys
//  * "frame_id"
//  * "regions"
//  * "frame_name"
//  could be addressed with UDSF's arguments. When the arguments are empty, this
//  stream function applies default key name.
//
// Stream Tuple.Data structure:
//  data.Map{
//    "frame_id": [frame id] (`data.Int`),
//    "regions" : [detected regions] (`[]data.Blob`)
//    "image":    [captured image] (`data.Blob`)
//  }
func (c *PredictTagsBatchStreamFuncCreator) CreateStreamFunction() interface{} {
	return createPredictTagsBatchUDSF
}

// TypeName returns type name.
func (c *PredictTagsBatchStreamFuncCreator) TypeName() string {
	return "scouter_predict_tags_batch_stream"
}

// CroppingAndPredictTagsBatchFuncCreator is a creator of cropping and predict
// tags batch UDF.
type CroppingAndPredictTagsBatchFuncCreator struct{}

// CreateFunction returns cropping and predict tags batch function. This
// function executes two tasks. First, cropping an image took by tagger
// parameters. Second, predicting tags and return regions with the tags. Tags
// information is set in a region.
//
// Usage:
//  `scouter_crop_and_predict_tags_batch([tagger_param], [regions], [image])`
//  [tagger_param]
//    * type: string
//    * a parameter name of "scouter_image_tagger_caffe" state
//  [region]
//    * type: []data.Blob
//    * detected regions created by detected function.
//    * these regions are detected from [image]
//  [image]
//    * type: []byte
//    * a captured image
//
// Return:
//  The function will return tagging regions, the type is `[]data.Blob`.
func (c *CroppingAndPredictTagsBatchFuncCreator) CreateFunction() interface{} {
	return croppingAndPredictTagsBatch
}

// TypeName returns type name.
func (c *CroppingAndPredictTagsBatchFuncCreator) TypeName() string {
	return "scouter_crop_and_predict_tags_batch"
}

func croppingAndPredictTagsBatch(ctx *core.Context, taggerParam string,
	regions data.Array, img []byte) (data.Array, error) {
	defer utils.LogElapseTime(ctx, "croppingAndPredictTagsBatch", time.Now())

	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	image := bridge.DeserializeMatVec3b(img)
	defer image.Delete()

	cans := make([]bridge.Candidate, len(regions))
	for i, r := range regions {
		regionByte, err := data.ToBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		cans[i] = regionPtr
	}

	defer func() {
		for _, c := range cans {
			c.Delete()
		}
	}()

	recognized := s.tagger.CropAndPredictTagsBatch(cans, image)

	defer func() {
		for _, r := range recognized {
			r.Delete()
		}
	}()

	recognizedCans := make(data.Array, len(recognized))
	for i, r := range recognized {
		recognizedCans[i] = data.Blob(r.Serialize())
	}

	return recognizedCans, nil
}

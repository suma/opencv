package recog

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"time"
)

// ACFDetectBatchFuncCreator is a creator of cropping region from frame UDF.
type RegionCropFuncCreator struct{}

func crop(ctx *core.Context, taggerParam string, region []byte, image []byte) (
	[]byte, error) {
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	r := bridge.DeserializeCandidate(region)
	defer r.Delete()

	img := bridge.DeserializeMatVec3b(image)
	defer img.Delete()

	cropped := s.tagger.Crop(r, img)
	defer cropped.Delete()
	return cropped.Serialize(), nil
}

// CreateFunction returns cropping region from frame function.
//
// Usage:
//  `scouter_crop_region([tagger_param], [region], [image])`
//  [tagger_param]
//    * type: string
//    * a parameter name of "scouter_image_tagger_caffe" state
//  [region]
//    * type: []byte
//    * a detected region created by detector UDF
//    * the region is detected from [image]
//  [image]
//    * type: []byte
//    * a captured image
//
// Return:
//  The function will return cropped imaged, the type is `[]byte`.
func (c *RegionCropFuncCreator) CreateFunction() interface{} {
	return crop
}

func (c *RegionCropFuncCreator) TypeName() string {
	return "scouter_crop_region"
}

type predictTagsBatchUDSF struct {
	predictTagsBatch func(
		[]bridge.Candidate, []bridge.MatVec3b) []bridge.Candidate
	frameIDName      string
	regionsName      string
	croppedImageName string
	detectCount      detectCounter
}

type detectCounter struct {
	sync.RWMutex
	count map[string]int
}

func (c *detectCounter) get(k string) (int, bool) {
	c.RLock()
	defer c.RUnlock()
	prev, ok := c.count[k]
	return prev, ok
}

func (c *detectCounter) put(k string, v int) {
	c.Lock()
	defer c.Unlock()
	c.count[k] = v
}

// Process streams tagged regions, which is serialized from
// `scouter::ObjectCandidate`. Tags information is set in a region.
//
// Stream Tuple.Data structure:
//  data.Map{
//    "frame_id": [frame ID] (`data.Int`),
//    "region":   [tagged region] (`data.Blob`),
//  }
func (sf *predictTagsBatchUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {
	frameID, err := t.Data.Get(sf.frameIDName)
	if err != nil {
		return err
	}
	frameIDStr, err := data.ToString(frameID)
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

	croppedImgsData, err := t.Data.Get(sf.croppedImageName)
	if err != nil {
		return err
	}
	croppedImgs, err := data.AsArray(croppedImgsData)
	if err != nil {
		return err
	}

	if len(regions) != len(croppedImgs) {
		return fmt.Errorf(
			"region size and cropped image size must same [region: %d, cropped image: %d",
			len(regions), len(croppedImgs))
	}

	if prevCount, ok := sf.detectCount.get(frameIDStr); ok {
		if prevCount > len(regions) {
			ctx.Log().Debug("prediction has already created")
			return nil
		}
	}
	sf.detectCount.put(frameIDStr, len(regions))

	candidates := []bridge.Candidate{}
	cropps := []bridge.MatVec3b{}
	defer func() {
		for _, c := range candidates {
			c.Delete()
		}
		for _, c := range cropps {
			c.Delete()
		}
	}()
	for i, r := range regions {
		rb, err := data.AsBlob(r)
		if err != nil {
			return err
		}
		candidates = append(candidates, bridge.DeserializeCandidate(rb))

		cb, err := data.AsBlob(croppedImgs[i])
		if err != nil {
			return err
		}
		cropps = append(cropps, bridge.DeserializeMatVec3b(cb))
	}

	recognized := sf.predictTagsBatch(candidates, cropps)
	defer func() {
		for _, r := range recognized {
			r.Delete()
		}
	}()

	traceCopyFlag := len(t.Trace) > 0
	for _, r := range recognized {
		now := time.Now()
		m := data.Map{
			"region_with_tagger": data.Blob(r.Serialize()),
			"frame_id":           frameID,
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

func (sf *predictTagsBatchUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createPredictTagsBatchUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	taggerParam string, stream string, frameIDName string, regionsName string,
	croppedImageName string) (udf.UDSF, error) {
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
		predictTagsBatch: s.tagger.PredictTagsBatch,
		frameIDName:      frameIDName,
		regionsName:      regionsName,
		croppedImageName: croppedImageName,
		detectCount:      detectCounter{count: map[string]int{}},
	}, nil
}

// PredictTagsBatchStreamFuncCreator is a creator of predicting tags UDSF.
type PredictTagsBatchStreamFuncCreator struct{}

// CreateStreamFunction returns predicting tags stream function. This stream
// function requires ID per frame to determine the regions detected from.
//
// Usage:
//  ```
//  scouter_predict_tags_batch_stream([tagger_param], [stream],
//                                    [frame_id_name], [regions_name])
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
//
// Input tuples are required to have following `data.Map` structures. The keys
//  * "frame_id"
//  * "regions"
//  could be addressed with UDSF's arguments. When the arguments are empty, this
//  stream function applies default key name.
//
// Stream Tuple.Data structure:
//  data.Map{
//    "frame_id": [frame id] (`data.Int`),
//    "regions" : [detected regions] (`[]data.Blob`)
//  }
func (c *PredictTagsBatchStreamFuncCreator) CreateStreamFunction() interface{} {
	return createPredictTagsBatchUDSF
}

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

	cans := []bridge.Candidate{}
	cropped := []bridge.MatVec3b{}
	for _, r := range regions {
		regionByte, err := data.AsBlob(r)
		if err != nil {
			return nil, err
		}
		regionPtr := bridge.DeserializeCandidate(regionByte)
		cans = append(cans, regionPtr)

		c := s.tagger.Crop(regionPtr, image)
		cropped = append(cropped, c)
	}

	defer func() {
		for _, c := range cans {
			c.Delete()
		}
		for _, c := range cropped {
			c.Delete()
		}
	}()

	recognized := s.tagger.PredictTagsBatch(cans, cropped)

	defer func() {
		for _, r := range recognized {
			r.Delete()
		}
	}()

	recognizedCans := data.Array{}
	for _, r := range recognized {
		recognizedCans = append(recognizedCans, data.Blob(r.Serialize()))
	}

	return recognizedCans, nil
}

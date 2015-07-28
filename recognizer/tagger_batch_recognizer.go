package recog

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"time"
)

type RegionCropFuncCreator struct{}

func crop(ctx *core.Context, taggerParam string, region data.Blob, image data.Blob) (data.Value, error) {
	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	regionByte, err := data.AsBlob(region)
	if err != nil {
		return nil, err
	}
	r := bridge.DeserializeCandidate(regionByte)
	defer r.Delete()

	imageByte, err := data.AsBlob(image)
	if err != nil {
		return nil, err
	}
	img := bridge.DeserializeMatVec3b(imageByte)
	defer img.Delete()

	cropped := s.tagger.Crop(r, img)
	defer cropped.Delete()
	return data.Blob(cropped.Serialize()), nil
}

func (c *RegionCropFuncCreator) CreateFunction() interface{} {
	return crop
}

func (c *RegionCropFuncCreator) TypeName() string {
	return "region_crop"
}

type predictTagsBatchUDSF struct {
	predictTagsBatch      func([]bridge.Candidate, []bridge.MatVec3b) []bridge.Candidate
	frameIDFieldName      string
	regionsFieldName      string
	croppedImageFieldName string
	detectCount           detectCounter
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

func (sf *predictTagsBatchUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	frameID, err := t.Data.Get(sf.frameIDFieldName)
	if err != nil {
		return err
	}
	frameIDStr, err := data.ToString(frameID)
	if err != nil {
		return err
	}

	regionsData, err := t.Data.Get(sf.regionsFieldName)
	if err != nil {
		return err
	}
	regions, err := data.AsArray(regionsData)
	if err != nil {
		return err
	}

	croppedImgsData, err := t.Data.Get(sf.croppedImageFieldName)
	if err != nil {
		return err
	}
	croppedImgs, err := data.AsArray(croppedImgsData)
	if err != nil {
		return err
	}

	if len(regions) != len(croppedImgs) {
		return fmt.Errorf("region size and cropped image size must same [region: %d, cropped image: %d",
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

func createPredictTagsBatchUDSF(ctx *core.Context, decl udf.UDSFDeclarer, taggerParam string,
	stream string, frameIDFieldName string, regionsFieldName string,
	croppedImageFieldName string) (udf.UDSF, error) {
	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "predict_tags_batch_stream",
	}); err != nil {
		return nil, err
	}

	s, err := lookupImageTaggerCaffeParamState(ctx, taggerParam)
	if err != nil {
		return nil, err
	}

	return &predictTagsBatchUDSF{
		predictTagsBatch:      s.tagger.PredictTagsBatch,
		frameIDFieldName:      frameIDFieldName,
		regionsFieldName:      regionsFieldName,
		croppedImageFieldName: croppedImageFieldName,
		detectCount:           detectCounter{count: map[string]int{}},
	}, nil
}

type PredictTagsBatchStreamFuncCreator struct{}

func (c *PredictTagsBatchStreamFuncCreator) CreateStreamFunction() interface{} {
	return createPredictTagsBatchUDSF
}

func (c *PredictTagsBatchStreamFuncCreator) TypeName() string {
	return "predict_tags_batch_stream"
}

func lookupImageTaggerCaffeParamState(ctx *core.Context, taggerParam string) (*ImageTaggerCaffeParamState, error) {
	st, err := ctx.SharedStates.Get(taggerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ImageTaggerCaffeParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to image_tagger_caffe_param.state", taggerParam)
}

package recog

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
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

type predictTagsBatchUDSF struct {
	predictTagsBatch      func([]bridge.Candidate, []bridge.MatVec3b) []bridge.Candidate
	frameIdFieldName      string
	regionsFieldName      string
	croppedImageFieldName string
}

func (sf *predictTagsBatchUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	frameId, err := t.Data.Get(sf.frameIdFieldName)
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

	candidates := []bridge.Candidate{}
	cropps := []bridge.MatVec3b{}
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
	for _, r := range recognized {
		now := time.Now()
		m := data.Map{
			"region_with_tagger": data.Blob(r.Serialize()),
			"frame_id":           frameId,
		}
		tu := &core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: t.ProcTimestamp,
			Trace:         make([]core.TraceEvent, 0),
		}
		w.Write(ctx, tu)
	}
	return nil
}

func (sf *predictTagsBatchUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func CreatePredictTagsBatchUDSF(ctx *core.Context, decl udf.UDSFDeclarer, taggerParam string,
	stream string, frameIdFieldName string, regionsFieldName string,
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
		frameIdFieldName:      frameIdFieldName,
		regionsFieldName:      regionsFieldName,
		croppedImageFieldName: croppedImageFieldName,
	}, nil
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

package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

func ACFDetectFunc(ctx *core.Context, detectParam string, frame data.Map) (data.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	img, err := lookupFrameData(frame)
	if err != nil {
		return nil, err
	}
	offsetX, offsetY, err := loopupOffsets(frame)
	if err != nil {
		return nil, err
	}
	imgP := bridge.DeserializeMatVec3b(img)
	defer imgP.Delete()
	candidates := s.d.ACFDetect(imgP, offsetX, offsetY)
	detected := data.Array{}
	for _, candidate := range candidates {
		detected = append(detected, data.Blob(candidate.Serialize()))
	}
	frame["detect"] = detected
	return frame, nil
}

type acfDetectUDSF struct {
	acfDetect          func(bridge.MatVec3b, int, int) []bridge.Candidate
	frameIdFieldName   string
	frameDataFieldName string
}

func (sf *acfDetectUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	frameId, err := t.Data.Get(sf.frameIdFieldName)
	if err != nil {
		return err
	}

	frame, err := t.Data.Get(sf.frameDataFieldName)
	if err != nil {
		return err
	}
	frameMeta, err := data.AsMap(frame)
	if err != nil {
		return err
	}

	img, err := lookupFrameData(frameMeta)
	if err != nil {
		return err
	}
	offsetX, offsetY, err := loopupOffsets(frameMeta)
	if err != nil {
		return err
	}

	imgP := bridge.DeserializeMatVec3b(img)
	defer imgP.Delete()
	candidates := sf.acfDetect(imgP, offsetX, offsetY)
	for _, c := range candidates {
		m := data.Map{
			"region":   data.Blob(c.Serialize()),
			"frame_id": frameId,
		}
		t.Data = m
		w.Write(ctx, t)
	}
	return nil
}

func (d *acfDetectUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func CreateACFDetectUDSF(ctx *core.Context, decl udf.UDSFDeclarer, detectParam string,
	stream string, frameIdFieldName string, frameDataFieldName string) (udf.UDSF, error) {
	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "acf_detector_stream",
	}); err != nil {
		return nil, err
	}

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	return &acfDetectUDSF{
		acfDetect:          s.d.ACFDetect,
		frameIdFieldName:   frameIdFieldName,
		frameDataFieldName: frameDataFieldName,
	}, nil
}

func FilterByMaskFunc(ctx *core.Context, detectParam string, region data.Blob) (data.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	r, err := data.AsBlob(region)
	if err != nil {
		return nil, err
	}

	regionPtr := bridge.DeserializeCandidate(r)
	defer regionPtr.Delete()
	masked := s.d.FilterByMask(regionPtr)
	return data.Bool(!masked), nil
}

func EstimateHeightFunc(ctx *core.Context, detectParam string, frame data.Map, region data.Blob) (data.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	offsetX, offsetY, err := loopupOffsets(frame)
	if err != nil {
		return nil, err
	}

	r, err := data.AsBlob(region)
	if err != nil {
		return nil, err
	}

	regionPtr := bridge.DeserializeCandidate(r)
	defer regionPtr.Delete()
	s.d.EstimateHeight(&regionPtr, offsetX, offsetY)

	return data.Blob(regionPtr.Serialize()), nil
}

func DrawDetectionResultFunc(ctx *core.Context, frame data.Blob, regions data.Array) (data.Value, error) {
	b, err := data.AsBlob(frame)
	if err != nil {
		return nil, err
	}
	img := bridge.DeserializeMatVec3b(b)

	canObjs := []bridge.Candidate{}
	for _, c := range regions {
		b, err := data.AsBlob(c)
		if err != nil {
			return nil, err
		}
		canObjs = append(canObjs, bridge.DeserializeCandidate(b))
	}

	ret := bridge.DrawDetectionResult(img, canObjs)
	return data.Blob(ret.Serialize()), nil
}

func lookupACFDetectParamState(ctx *core.Context, detectParam string) (*ACFDetectionParamState, error) {
	st, err := ctx.GetSharedState(detectParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*ACFDetectionParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to acf_detection_param.state", detectParam)
}

func lookupFrameData(frame data.Map) ([]byte, error) {
	img, err := frame.Get("projected_img")
	if err != nil {
		return []byte{}, err
	}
	image, err := data.AsBlob(img)
	if err != nil {
		return []byte{}, err
	}

	return image, nil
}

func loopupOffsets(frame data.Map) (int, int, error) {
	ox, err := frame.Get("offset_x")
	if err != nil {
		return 0, 0, err
	}
	offsetX, err := data.AsInt(ox)
	if err != nil {
		return 0, 0, err
	}

	oy, err := frame.Get("offset_y")
	if err != nil {
		return 0, 0, err
	}
	offsetY, err := data.AsInt(oy)
	if err != nil {
		return 0, 0, err
	}

	return int(offsetX), int(offsetY), nil
}

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

	_, img, err := lookupFrameData(frame)
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
	id         int64
	candidates []bridge.Candidate
}

func (sf *acfDetectUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	for _, c := range sf.candidates {
		m := data.Map{
			"frame_id": data.Int(sf.id),
			"region":   data.Blob(c.Serialize()),
		}
		t.Data = m
		w.Write(ctx, t)
	}
	return nil
}

func (d *acfDetectUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func CreateACFDetectUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	detectParam string, frame data.Map) (udf.UDSF, error) {
	// TODO declarer input

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	id, img, err := lookupFrameData(frame)
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

	return &acfDetectUDSF{
		id:         id,
		candidates: candidates,
	}, nil
}

func FilterByMaskFunc(ctx *core.Context, detectParam string, frame data.Map) (data.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	candidates, err := frame.Get("detect")
	if err != nil {
		return nil, err
	}
	cans, err := data.AsArray(candidates)
	if err != nil {
		return nil, err
	}

	canObjs := []bridge.Candidate{}
	for _, c := range cans {
		b, err := data.AsBlob(c)
		if err != nil {
			return nil, err // TODO return is OK?
		}
		canObjs = append(canObjs, bridge.DeserializeCandidate(b))
	}

	filteredCans := s.d.FilterByMask(canObjs)
	filtered := data.Array{}
	for _, fc := range filteredCans {
		filtered = append(filtered, data.Blob(fc.Serialize()))
		fc.Delete() // TODO use defer
	}

	frame["detect"] = filtered // TODO overwrite is OK?
	return frame, nil
}

func EstimateHeightFunc(ctx *core.Context, detectParam string, frame data.Map) (data.Value, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	offsetX, offsetY, err := loopupOffsets(frame)
	if err != nil {
		return nil, err
	}

	candidates, err := frame.Get("detect")
	if err != nil {
		return nil, err
	}
	cans, err := data.AsArray(candidates)
	if err != nil {
		return nil, err
	}

	canObjs := []bridge.Candidate{}
	for _, c := range cans {
		b, err := data.AsBlob(c)
		if err != nil {
			return nil, err // TODO return is OK?
		}
		canObjs = append(canObjs, bridge.DeserializeCandidate(b))
	}

	estimatedCans := s.d.EstimateHeight(canObjs, offsetX, offsetY)
	estimated := data.Array{}
	for _, ec := range estimatedCans {
		estimated = append(estimated, data.Blob(ec.Serialize()))
	}

	frame["detect"] = estimated // TODO overwrite is OK?
	return frame, nil
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

func lookupFrameData(frame data.Map) (int64, []byte, error) {
	id, err := frame.Get("frame_id")
	if err != nil {
		return 0, []byte{}, err
	}
	frameId, err := data.AsInt(id)
	if err != nil {
		return 0, []byte{}, err
	}

	img, err := frame.Get("projected_img")
	if err != nil {
		return 0, []byte{}, err
	}
	image, err := data.AsBlob(img)
	if err != nil {
		return 0, []byte{}, err
	}

	return frameId, image, nil
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

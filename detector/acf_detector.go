package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

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
	defer func() {
		for _, c := range candidates {
			c.Delete()
		}
	}()
	for _, c := range candidates {
		now := time.Now()
		m := data.Map{
			"region":   data.Blob(c.Serialize()),
			"frame_id": frameId,
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

func (sf *acfDetectUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createACFDetectUDSF(ctx *core.Context, decl udf.UDSFDeclarer, detectParam string,
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

type DetectRegionStreamFuncCreator struct{}

func (c *DetectRegionStreamFuncCreator) CreateStreamFunction() interface{} {
	return createACFDetectUDSF
}

func (c *DetectRegionStreamFuncCreator) TypeName() string {
	return "acf_detector_stream"
}

type FilterByMaskFuncCreator struct{}

func filterByMask(ctx *core.Context, detectParam string, region data.Blob) (bool, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return false, err
	}

	r, err := data.AsBlob(region)
	if err != nil {
		return false, err
	}

	regionPtr := bridge.DeserializeCandidate(r)
	defer regionPtr.Delete()
	masked := s.d.FilterByMask(regionPtr)

	return !masked, nil
}

func (c *FilterByMaskFuncCreator) CreateFunction() interface{} {
	return filterByMask
}

func (c *FilterByMaskFuncCreator) TypeName() string {
	return "filter_by_mask"
}

type EstimateHeightFuncCreator struct{}

func estimateHeight(ctx *core.Context, detectParam string, frame data.Map, region data.Blob) (data.Value, error) {
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

func (c *EstimateHeightFuncCreator) CreateFunction() interface{} {
	return estimateHeight
}

func (c *EstimateHeightFuncCreator) TypeName() string {
	return "estimate_height"
}

type DrawDetectionResultFuncCreator struct{}

func drawDetectionResult(ctx *core.Context, frame data.Blob, regions data.Array) (data.Value, error) {
	b, err := data.AsBlob(frame)
	if err != nil {
		return nil, err
	}
	img := bridge.DeserializeMatVec3b(b)
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

	ret := bridge.DrawDetectionResult(img, canObjs)
	defer ret.Delete()
	return data.Blob(ret.Serialize()), nil
}

func (c *DrawDetectionResultFuncCreator) CreateFunction() interface{} {
	return drawDetectionResult
}

func (c *DrawDetectionResultFuncCreator) TypeName() string {
	return "draw_detection_result"
}

func lookupACFDetectParamState(ctx *core.Context, detectParam string) (*ACFDetectionParamState, error) {
	st, err := ctx.SharedStates.Get(detectParam)
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
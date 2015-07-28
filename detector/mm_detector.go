package detector

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type mmDetectUDSF struct {
	mmDetect           func(bridge.MatVec3b, int, int) []bridge.Candidate
	frameIdFieldName   string
	frameDataFieldName string
}

func (sf *mmDetectUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
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
	candidates := sf.mmDetect(imgP, offsetX, offsetY)
	defer func() {
		for _, c := range candidates {
			c.Delete()
		}
	}()

	traceCopyFlag := len(t.Trace) > 0
	for _, c := range candidates {
		now := time.Now()
		m := data.Map{
			"region":   data.Blob(c.Serialize()),
			"frame_id": frameId,
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

func (sf *mmDetectUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createMMDetectUDSF(ctx *core.Context, decl udf.UDSFDeclarer, detectParam string,
	stream string, frameIdFieldName string, frameDataFieldName string) (udf.UDSF, error) {
	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "multi_model_detector_stream",
	}); err != nil {
		return nil, err
	}

	s, err := lookupMMDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	return &mmDetectUDSF{
		mmDetect:           s.d.MMDetect,
		frameIdFieldName:   frameIdFieldName,
		frameDataFieldName: frameDataFieldName,
	}, nil
}

type MMDetectRegionStreamFuncCreator struct{}

func (c *MMDetectRegionStreamFuncCreator) CreateStreamFunction() interface{} {
	return createMMDetectUDSF
}

func (c *MMDetectRegionStreamFuncCreator) TypeName() string {
	return "multi_model_detector_stream"
}

type FilterByMaskMMFuncCreator struct{}

func filterByMaskMM(ctx *core.Context, detectParam string, region data.Blob) (bool, error) {
	s, err := lookupMMDetectParamState(ctx, detectParam)
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

func (c *FilterByMaskMMFuncCreator) CreateFunction() interface{} {
	return filterByMaskMM
}

func (c *FilterByMaskMMFuncCreator) TypeName() string {
	return "multi_model_filter_by_mask"
}

type EstimateHeightMMFuncCreator struct{}

func estimateHeightMM(ctx *core.Context, detectParam string, frame data.Map, region data.Blob) (data.Value, error) {
	s, err := lookupMMDetectParamState(ctx, detectParam)
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

func (c *EstimateHeightMMFuncCreator) CreateFunction() interface{} {
	return estimateHeightMM
}

func (c *EstimateHeightMMFuncCreator) TypeName() string {
	return "multi_model_estimate_height"
}

func lookupMMDetectParamState(ctx *core.Context, detectParam string) (*MMDetectionParamState, error) {
	st, err := ctx.SharedStates.Get(detectParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*MMDetectionParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to mm_detection_param.state", detectParam)
}

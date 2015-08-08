package detector

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type mmDetectUDSF struct {
	mmDetect    func(bridge.MatVec3b, int, int) []bridge.Candidate
	frameIDName string
	frameName   string
}

// Process streams detected regions. which is serialized from
// `scouter::ObjectCandidate`.
//
//  data.Map{
//    "frame_id":      [frame ID] (`data.Int`),
//    "regions_count": [size of regions created from frame] (`data.Int`),
//    "region":        [detected region] (`data.Blob`),
//  }
func (sf *mmDetectUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	frameId, err := t.Data.Get(sf.frameIDName)
	if err != nil {
		return err
	}

	frame, err := t.Data.Get(sf.frameName)
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
	offsetX, offsetY, err := lookupOffsets(frameMeta)
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
			"frame_id":      frameId,
			"recions_count": data.Int(len(candidates)),
			"region":        data.Blob(c.Serialize()),
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
	stream string, frameIDName string, frameName string) (udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "scouter_mm_detector_stream",
	}); err != nil {
		return nil, err
	}

	s, err := lookupMMDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	if frameIDName == "" {
		frameIDName = "frame_id"
	}
	if frameName == "" {
		frameName = "frame"
	}

	return &mmDetectUDSF{
		mmDetect:    s.d.MMDetect,
		frameIDName: frameIDName,
		frameName:   frameName,
	}, nil
}

// MMDetectRegionStreamFuncCreator is a creator of Multi Model detector UDSF.
type MMDetectRegionStreamFuncCreator struct{}

// CreateStreamFunction returns Multi Model Detection stream function. This
// stream function requires ID per frame to determine the regions detected from.
//
// Usage:
//  ```
//  scouter_mm_detector_stream([detect_param], [stream],
//                             [frame_id_name], [frame_name])
//  ```
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_mm_detection_param" state
//  [stream]
//    * type: string
//    * a input stream name, see following.
//  [frame_id_name]
//    * type: string
//    * a field name of frame ID
//    * if empty then applied "frame_id"
//  [frame_name]
//    * type: string
//    * a field name of frame
//    * if empty then applied "frame"
//
// Input tuples are required to have following `data.Map` structure. The two
// keys
//   * "frame_id"
//   * "frame"
// could be addressed with UDS's arguments. When the arguments are empty,
// this stream function applies default key name.
//
//
//  data.Map{
//    "frame_id": [frame id] (`data.Int`),
//    "frame"   : data.Map{
//      "projected_img": [image binary] (`data.Blob`),
//      "offset_x":      [frame offset x] (`data.Int`),
//      "offset_y":      [frame offset y] (`data.Int`),
//    }
//  }
func (c *MMDetectRegionStreamFuncCreator) CreateStreamFunction() interface{} {
	return createMMDetectUDSF
}

func (c *MMDetectRegionStreamFuncCreator) TypeName() string {
	return "scouter_mm_detector_stream"
}

// FilterByMaskMMFuncCreator is a creator of filtering by mask UDF.
type FilterByMaskMMFuncCreator struct{}

func filterByMaskMM(ctx *core.Context, detectParam string, region []byte) (
	bool, error) {
	s, err := lookupMMDetectParamState(ctx, detectParam)
	if err != nil {
		return false, err
	}

	regionPtr := bridge.DeserializeCandidate(region)
	defer regionPtr.Delete()
	masked := s.d.FilterByMask(regionPtr)

	return !masked, nil
}

// CreateFunction creates a filter by mask for Multi Model detection.
//
// Usage:
//  `scouter_mm_filter_by_mask([detect_param], [region])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_mm_detection_param" state
//  [region]
//    * type: []byte
//    * detected region, which are applied Multi Model detection.
//
// Returns:
//  The function will return the region is filtered or not, the type is `bool`.
func (c *FilterByMaskMMFuncCreator) CreateFunction() interface{} {
	return filterByMaskMM
}

func (c *FilterByMaskMMFuncCreator) TypeName() string {
	return "scouter_mm_filter_by_mask"
}

// EstimateHeightMMFuncCreator is creator of height estimator UDF.
type EstimateHeightMMFuncCreator struct{}

func estimateHeightMM(ctx *core.Context, detectParam string, frame data.Map,
	region []byte) ([]byte, error) {
	s, err := lookupMMDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	offsetX, offsetY, err := lookupOffsets(frame)
	if err != nil {
		return nil, err
	}

	regionPtr := bridge.DeserializeCandidate(region)
	defer regionPtr.Delete()
	s.d.EstimateHeight(&regionPtr, offsetX, offsetY)
	return regionPtr.Serialize(), nil
}

// CreateFunction creates a estimate height function for Multi Model detection.
//
// Usage:
//  `scouter_mm_estimate_height([detect_param], [frame], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_mm_detection_param" state
//  [frame]
//    * type: data.Map
//    * captured frame which are applied `scouter_frame_applier` UDF. The
//      frame's map structure is required following structure.
//      data.Map{
//        "offset_x"  : [frame offset x] (`data.Int`)
//        "offset_y"  : [frame offset y] (`data.Int`)
//      }
//  [region]
//    * type: []byte
//    * detected region, which are applied Multi Model detection.
//    * the region is detected from [frame]
//
// Return:
//   The function will return an estimate region, the type is `[]byte`.
func (c *EstimateHeightMMFuncCreator) CreateFunction() interface{} {
	return estimateHeightMM
}

func (c *EstimateHeightMMFuncCreator) TypeName() string {
	return "scouter_mm_estimate_height"
}

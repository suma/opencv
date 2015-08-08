package detector

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type acfDetectUDSF struct {
	acfDetect   func(bridge.MatVec3b, int, int) []bridge.Candidate
	frameIDName string
	frameName   string
}

// Process streams detected regions. which is serialized from
// `scouter::ObjectCandidate`.
//
//  data.Map{
//    "region"  : [detected region] (`data.Blob`),
//    "frame_id": [frame ID] (`data.Int`),
//  }
func (sf *acfDetectUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {

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
	candidates := sf.acfDetect(imgP, offsetX, offsetY)
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

func (sf *acfDetectUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createACFDetectUDSF(ctx *core.Context, decl udf.UDSFDeclarer, detectParam string,
	stream string, frameIDName string, frameName string) (udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "scouter_acf_detector_stream",
	}); err != nil {
		return nil, err
	}

	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return nil, err
	}

	if frameIDName == "" {
		frameIDName = "frame_id"
	}
	if frameName == "" {
		frameName = "frame"
	}

	return &acfDetectUDSF{
		acfDetect:   s.d.ACFDetect,
		frameIDName: frameIDName,
		frameName:   frameName,
	}, nil
}

// DetectRegionStreamFuncCreator is a creator of ACF detector UDSF.
type DetectRegionStreamFuncCreator struct{}

// CreateStreamFunction returns ACF Detection stream function. This stream
// function requires ID per frame to determine the regions detected from.
//
// Usage:
//  ```
//  scouter_acf_detector_stream([detect_param], [stream],
//                              [frame_id_name], [frame_name])
//  ```
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_detection_param" state
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
func (c *DetectRegionStreamFuncCreator) CreateStreamFunction() interface{} {
	return createACFDetectUDSF
}

func (c *DetectRegionStreamFuncCreator) TypeName() string {
	return "scouter_acf_detector_stream"
}

// FilterByMaskFuncCreator is a creator of filtering by mask UDF.
type FilterByMaskFuncCreator struct{}

func filterByMask(ctx *core.Context, detectParam string, region []byte) (bool, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
	if err != nil {
		return false, err
	}

	regionPtr := bridge.DeserializeCandidate(region)
	defer regionPtr.Delete()
	masked := s.d.FilterByMask(regionPtr)

	return !masked, nil
}

// CreateFunction creates a filter by mask for ACF detection.
//
// Usage:
//  `scouter_filter_by_mask([detect_param], [region])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_detection_param" state
//  [region]
//    * type: []byte
//    * detected region, which are applied ACF detection.
//
// Returns:
//  The function will return the region is filtered or not, the type is `bool`.
func (c *FilterByMaskFuncCreator) CreateFunction() interface{} {
	return filterByMask
}

func (c *FilterByMaskFuncCreator) TypeName() string {
	return "scouter_filter_by_mask"
}

// EstimateHeightFuncCreator is creator of height estimator UDF.
type EstimateHeightFuncCreator struct{}

func estimateHeight(ctx *core.Context, detectParam string, frame data.Map,
	region []byte) ([]byte, error) {
	s, err := lookupACFDetectParamState(ctx, detectParam)
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

// CreateFunction creates a estimate height function for ACF detection.
//
// Usage:
//  `scouter_estimate_height([detect_param], [frame], [regions])`
//  [detect_param]
//    * type: string
//    * a parameter name of "scouter_detection_param" state
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
//    * detected region, which are applied ACF detection.
//    * the region is detected from [frame]
//
// Return:
//   The function will return an estimate region, the type is `[]byte`.
func (c *EstimateHeightFuncCreator) CreateFunction() interface{} {
	return estimateHeight
}

func (c *EstimateHeightFuncCreator) TypeName() string {
	return "scouter_estimate_height"
}

// DrawDetectionResultFuncCreator is a creator of drawing regions on a frame UDF.
type DrawDetectionResultFuncCreator struct{}

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

	ret := bridge.DrawDetectionResult(img, canObjs)
	defer ret.Delete()
	return ret.Serialize(), nil
}

// CreateFunction creates a drawing regions on a frame function.
//
// Usage:
//  `scouter_draw_regions([frame], [regions])`
//  [frame]
//    * type: []byte
//    * captured frame, which is serialized from `cv::Mat_<cv::Vec3b>`.
//  [regions]
//    * type: []data.Blob
//    * detected regions, which are applied detector UDF/UDSF
//    * these regions are detected from [frame]
//
// Return:
//  The function will return an image data serialized from `cv::Mat_<cv::Vec3b>`,
//  the type is `[]byte`
func (c *DrawDetectionResultFuncCreator) CreateFunction() interface{} {
	return drawDetectionResult
}

func (c *DrawDetectionResultFuncCreator) TypeName() string {
	return "scouter_draw_regions"
}

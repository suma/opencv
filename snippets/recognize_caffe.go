package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type RecognizeCaffeConfig struct {
	PlayerFlag bool
}

type RecognizeCaffe struct {
	Config  RecognizeCaffeConfig
	taggers bridge.ImageTaggerCaffe
}

type FrameInfo struct {
	index int
	fr    []byte
	dr    []byte
}

func (rc *RecognizeCaffe) Init(ctx *core.Context) error {
	// TODO create configuration
	taggers := bridge.ImageTaggerCaffe_New(bridge.RecognizeConfigTaggers{})
	rc.taggers = taggers
	return nil
}

func (rc *RecognizeCaffe) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	fi, err := getFrameInfo(t)
	if err != nil {
		return err
	}

	rc.governor(fi)
	rc.recognize(fi, t)

	w.Write(ctx, t)
	return nil
}

func getFrameInfo(t *tuple.Tuple) (FrameInfo, error) {
	f, err := t.Data.Get("frame")
	if err != nil {
		return FrameInfo{}, fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return FrameInfo{}, fmt.Errorf("frame data must be byte array type")
	}

	d, err := t.Data.Get("detection_result")
	if err != nil {
		return FrameInfo{}, fmt.Errorf("cannot get detection result")
	}
	detectionResult, err := d.AsBlob()
	if err != nil {
		return FrameInfo{}, fmt.Errorf("detection result data must be byte array type")
	}

	return FrameInfo{
		fr: frame,
		dr: detectionResult}, nil
}

func (rc *RecognizeCaffe) governor(fi FrameInfo) {
	// join where meta.time is equal
}

func (rc *RecognizeCaffe) recognize(fi FrameInfo, t *tuple.Tuple) {
	fr := bridge.DeserializeFrame(fi.fr)
	dr := bridge.DeserializeDetectionResult(fi.dr)

	recogDr := rc.taggers.Recognize(fr, dr)
	t.Data["recognize_detection_result"] = tuple.Blob(recogDr.Serialize())

	if rc.Config.PlayerFlag {
		drwResult := bridge.RecognizeDrawResult(fr, recogDr)
		fmt.Println(drwResult)
		// TODO convert to map[string]
		//t.Data["recognize_draw_result"] = tuple.Blob(drwResult)
	}

	fr.Delete()
	dr.Delete() // TODO user defer
}

func (rc *RecognizeCaffe) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (rc *RecognizeCaffe) OutputSchema(ss []*core.Schema) (*core.Schema, error) {
	return nil, nil
}

func (rc *RecognizeCaffe) Terminate(ctx *core.Context) error {
	rc.taggers.Delete()
	return nil
}

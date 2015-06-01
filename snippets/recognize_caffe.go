package snippets

import (
	"fmt"
	"io/ioutil"
	"os"
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type RecognizeCaffe struct {
	ConfigPath string
	Config     conf.RecognizeCaffeConfig
	taggers    bridge.ImageTaggerCaffe
}

type FrameInfo struct {
	index int
	fr    []byte
	dr    []byte
}

func (rc *RecognizeCaffe) Init(ctx *core.Context) error {
	config, err := conf.GetRecognizeCaffeSnippetConfig(rc.ConfigPath)
	if err != nil {
		return err
	}
	rc.Config = config
	taggers := bridge.NewImageTaggerCaffe(config.ConfigTaggers)
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
	defer fr.Delete()
	dr := bridge.DeserializeDetectionResult(fi.dr)
	defer dr.Delete()

	recogDr := rc.taggers.Recognize(fr, dr)
	t.Data["recognize_detection_result"] = tuple.Blob(recogDr.Serialize())

	if rc.Config.PlayerFlag {
		drwResults := bridge.RecognizeDrawResult(fr, recogDr)
		for k, v := range drwResults {
			defer v.Delete()
			s := time.Now().UnixNano() / int64(time.Millisecond)
			t.Data["recognize_draw_result_"+k] = tuple.Blob(v.ToJpegData(rc.Config.JpegQuality))
			// following is debug for scouter recognize caffe
			ioutil.WriteFile(fmt.Sprintf("./recog_%v_%v.jpg", k, fmt.Sprint(s)),
				v.ToJpegData(50), os.ModePerm)
		}
	}
}

func (rc *RecognizeCaffe) Terminate(ctx *core.Context) error {
	rc.taggers.Delete()
	return nil
}

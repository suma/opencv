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

type DetectSimple struct {
	ConfigPath string
	Config     conf.DetectSimpleConfig
	detector   bridge.Detector
	lastFrame  *tuple.Tuple
}

func (d *DetectSimple) Init(ctx *core.Context) error {
	detectConfig, err := conf.GetDetectSimpleSnippetConfig(d.ConfigPath)
	if err != nil {
		return err
	}
	d.Config = detectConfig
	d.detector = bridge.NewDetector(detectConfig.DetectorConfig)
	d.lastFrame = nil
	return nil
}

func (d *DetectSimple) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	switch t.InputName {
	case "frame":
		d.lastFrame = t

	case "tick":
		if d.lastFrame == nil {
			return nil
		}
		// following process should be run in thread safe,
		// but following code is not thread safe.
		// tick interval is very longer than frame rate,
		// so detect process is implemented in simple copy strategy.
		frame := d.lastFrame.Copy()
		d.lastFrame = nil
		err := detect(d, frame, t.Timestamp)
		if err != nil {
			return err
		}

		w.Write(ctx, frame)
	}
	return nil
}

func detect(d *DetectSimple, t *tuple.Tuple, timestamp time.Time) error {
	f, err := getFrame(t)
	if err != nil {
		return err
	}

	fPointer := bridge.DeserializeFrame(f)
	defer fPointer.Delete()
	s := time.Now().UnixNano() / int64(time.Millisecond)
	drPointer := d.detector.Detect(fPointer)

	t.Data["detection_result"] = tuple.Blob(drPointer.Serialize())
	t.Data["detection_time"] = tuple.Timestamp(timestamp)

	if d.Config.PlayerFlag {
		ms := time.Now().UnixNano()/int64(time.Millisecond) - s
		drw := bridge.DetectDrawResult(fPointer, drPointer, ms)
		defer drw.Delete()
		t.Data["detection_draw_result"] = tuple.Blob(drw.ToJpegData(d.Config.JpegQuality))
		// following is debug for scouter detection
		ioutil.WriteFile(fmt.Sprintf("./detect_%v.jpg", fmt.Sprint(s)),
			drw.ToJpegData(50), os.ModePerm)
	}
	return nil
}

func getFrame(t *tuple.Tuple) ([]byte, error) {
	f, err := t.Data.Get("frame")
	if err != nil {
		return []byte{}, fmt.Errorf("cannot get frame data")
	}
	frame, err := f.AsBlob()
	if err != nil {
		return []byte{}, fmt.Errorf("frame data must be byte array type")
	}
	return frame, nil
}

func (d *DetectSimple) Terminate(ctx *core.Context) error {
	d.detector.Delete()
	return nil
}

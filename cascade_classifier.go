package opencv

import (
	"fmt"
	"pfi/sensorbee/opencv/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

var (
	configFilePath = data.MustCompilePath("file")
)

// NewCascadeClassifier returns cascadeClassifier state.
//
// file: cascade configuration file path for detection.
// e.g. "haarcascade_frontalface_default.xml".
func NewCascadeClassifier(ctx *core.Context, params data.Map) (core.SharedState,
	error) {
	var filePath string
	if fp, err := params.Get(configFilePath); err != nil {
		return nil, err
	} else if filePath, err = data.AsString(fp); err != nil {
		return nil, err
	}

	cc := bridge.NewCascadeClassifier()
	if !cc.Load(filePath) {
		return nil, fmt.Errorf("cannot load the file '%v'", filePath)
	}

	return &cascadeClassifier{
		classifier: cc,
	}, nil
}

type cascadeClassifier struct {
	classifier bridge.CascadeClassifier
}

func (c *cascadeClassifier) Terminate(ctx *core.Context) error {
	c.classifier.Delete()
	return nil
}

func lookupCascadeClassifier(ctx *core.Context, name string) (*cascadeClassifier,
	error) {
	st, err := ctx.SharedStates.Get(name)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*cascadeClassifier); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be canverted to cascade_classifier.state",
		name)
}

// DetectMultiScale classifies and detect image.
//
// classifierName: cascadeClassifier state name.
//
// img: target image as RawData map structure.
func DetectMultiScale(ctx *core.Context, classifierName string, img data.Map) (
	data.Array, error) {
	raw, err := ConvertMapToRawData(img)
	if err != nil {
		return nil, err
	}
	mat := raw.ToMatVec3b()
	defer mat.Delete()

	classifier, err := lookupCascadeClassifier(ctx, classifierName)
	if err != nil {
		return nil, err
	}
	rects := classifier.classifier.DetectMultiScale(mat)
	ret := make(data.Array, len(rects))
	for i, r := range rects {
		rect := data.Map{
			"x":      data.Int(r.X),
			"y":      data.Int(r.Y),
			"width":  data.Int(r.Width),
			"height": data.Int(r.Height),
		}
		ret[i] = rect
	}
	return ret, nil
}

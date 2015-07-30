package writer

import (
	"fmt"
	"os"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type VideoWiterCreator struct{}

func (c *VideoWiterCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	fn, err := params.Get("file_name")
	if err != nil {
		return nil, err
	}
	name, err := data.ToString(fn)
	if err != nil {
		return nil, err
	}
	name += ".avi"

	_, err = os.Stat(name)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%v has already been exist, cannot create video writer",
			name)
	}

	fps, err := params.Get("fps")
	if err != nil {
		fps = data.Float(1.0)
	}
	fpsRate, err := data.AsFloat(fps)
	if err != nil {
		return nil, err
	}

	w, err := params.Get("width")
	if err != nil {
		w = data.Int(480)
	}
	width, err := data.AsInt(w)
	if err != nil {
		return nil, err
	}

	h, err := params.Get("height")
	if err != nil {
		h = data.Int(320)
	}
	height, err := data.AsInt(h)
	if err != nil {
		return nil, err
	}

	vw := bridge.NewVideoWriter()
	vw.Open(name, fpsRate, int(width), int(height))
	if !vw.IsOpened() {
		return nil, fmt.Errorf("cannot video writer open: %v", name)
	}

	s := &videoWriterSink{}
	s.vw = vw
	return s, nil
}

func (c *VideoWiterCreator) TypeName() string {
	return "video_writer"
}

type videoWriterSink struct {
	vw bridge.VideoWriter
}

func (s *videoWriterSink) Write(ctx *core.Context, t *core.Tuple) error {
	img, err := t.Data.Get("img")
	if err != nil {
		return err
	}
	imgByte, err := data.AsBlob(img)
	if err != nil {
		return err
	}
	imgp := bridge.DeserializeMatVec3b(imgByte)
	defer imgp.Delete()

	s.vw.Write(imgp)
	return nil
}

func (s *videoWriterSink) Close(ctx *core.Context) error {
	s.vw.Delete()
	return nil
}

package mjpegserv

import (
	"fmt"
	"io/ioutil"
	"os"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type DebugJPEGSink struct {
	outputDir   string
	jpegQuality int
	detectCount map[string]int64
}

func (s *DebugJPEGSink) Write(ctx *core.Context, t *core.Tuple) error {
	count, err := t.Data.Get("region_count")
	if err != nil {
		return err
	}
	countInt, err := data.ToInt(count)
	if err != nil {
		return err
	}

	name, err := t.Data.Get("name")
	if err != nil {
		return err
	}
	nameStr, err := data.ToString(name)
	if err != nil {
		return err
	}

	if prevCount, ok := s.detectCount[nameStr]; ok {
		if prevCount > countInt {
			ctx.Log().Debug("JPEG has already created")
			return nil
		}
	}
	s.detectCount[nameStr] = countInt

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

	fileName := fmt.Sprintf("%v/%v.jpg", s.outputDir, nameStr)
	ioutil.WriteFile(fileName, imgp.ToJpegData(s.jpegQuality), os.ModePerm)
	return nil
}

func (s *DebugJPEGSink) Close(ctx *core.Context) error {
	return nil
}

func (s *DebugJPEGSink) CreateSink(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Sink, error) {
	output, err := params.Get("output")
	if err != nil {
		output = data.String(".")
	}
	outputDir, err := data.AsString(output)
	if err != nil {
		return nil, err
	}

	quality, err := params.Get("quality")
	if err != nil {
		quality = data.Int(50)
	}
	q, err := data.AsInt(quality)
	if err != nil {
		return nil, err
	}

	s.outputDir = outputDir
	s.jpegQuality = int(q)
	s.detectCount = map[string]int64{}
	return s, nil
}
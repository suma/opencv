package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// JPEGWriterCreator is a creator of JPEG Writer.
type JPEGWriterCreator struct{}

var (
	outputPath  = data.MustCompilePath("output")
	qualityPath = data.MustCompilePath("quality")
	namePath    = data.MustCompilePath("name")
	imgPath     = data.MustCompilePath("img")
)

// CreateSink creates a JPEG output sink, which output converted JPEG from
// `cv::Mat_<cv::Vec3b>`.
//
// Usage of WITH parameters:
//  output:  Output directory, If empty then files are output to the current
//           directory. If the directory is not exist, this sink will make the
//           output directory. Returns an error when the sink could not make it.
//  quality: The quality of converting JPEG file, if empty then set 50.
//
// Example:
//  when a creation query is
//    `CREATE SINK jpeg_files TYPE scouter_jpeg_writer
//      WITH output='temp', quality=50`
//  then will create JPEG files and output to "./temp" directory.
func (c *JPEGWriterCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	output, err := params.Get(outputPath)
	if err != nil {
		output = data.String(".")
	}
	outputDir, err := data.AsString(output)
	if err != nil {
		return nil, err
	}

	if absPath, err := filepath.Abs(outputDir); err != nil {
		return nil, fmt.Errorf("invalid file path: %v", err.Error())
	} else {
		_, err = os.Stat(absPath)
		if os.IsNotExist(err) {
			os.MkdirAll(absPath, 0755)
		}
	}

	quality, err := params.Get(qualityPath)
	if err != nil {
		quality = data.Int(50)
	}
	q, err := data.AsInt(quality)
	if err != nil {
		return nil, err
	}

	s := &jpegWriterSink{}
	s.outputDir = outputDir
	s.jpegQuality = int(q)
	return s, nil
}

func (c *JPEGWriterCreator) TypeName() string {
	return "scouter_jpeg_writer"
}

type jpegWriterSink struct {
	outputDir   string
	jpegQuality int
}

// Write input JPEG files to the directory which is set `WITH` "output"
// parameter. Input tuple is required to have follow `data.Map`
//
//  data.Map{
//    "name": [output file name] (will be casted to string type)
//    "img" : [image binary data] (`data.Blob`)
//  }
//
// Example of insertion query:
//  ```
//  INSERT INTO jpeg_files SELECT ISTREAM
//    frame_data AS img,
//    frame_id AS name
//    FROM capturing_frame [RANGE 1 TUPLES];
//  ```
// then [frame_id].jpg will be created at the directory.
func (s *jpegWriterSink) Write(ctx *core.Context, t *core.Tuple) error {
	name, err := t.Data.Get(namePath)
	if err != nil {
		return err
	}
	nameStr, err := data.ToString(name)
	if err != nil {
		return err
	}

	img, err := t.Data.Get(imgPath)
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

func (s *jpegWriterSink) Close(ctx *core.Context) error {
	return nil
}

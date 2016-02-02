package opencv

import (
	"bytes"
	"image"
	"image/jpeg"
	"pfi/sensorbee/opencv/bridge"
	"pfi/sensorbee/sensorbee/data"
)

var (
	imagePath = data.MustCompilePath("image")
)

// RawData is represented of `cv::Mat_<cv::Vec3b>` structure.
type RawData struct {
	Width  int
	Height int
	Data   []byte
}

// ToRawData converts MatVec3b to RawData.
func ToRawData(m bridge.MatVec3b) RawData {
	w, h, data := m.ToRawData()
	return RawData{
		Width:  w,
		Height: h,
		Data:   data,
	}
}

// ToMatVec3b converts RawData to MatVec3b. Returned MatVec3b is required to
// delete after using.
func (r *RawData) ToMatVec3b() bridge.MatVec3b {
	return bridge.ToMatVec3b(r.Width, r.Height, r.Data)
}

func toRawMap(m *bridge.MatVec3b) data.Map {
	r := ToRawData(*m)
	return data.Map{
		"format": data.String("cvmat"), // = cv::Mat_<cv::Vec3b>
		"width":  data.Int(r.Width),
		"height": data.Int(r.Height),
		"image":  data.Blob(r.Data),
	}
}

// ConvertMapToRawData returns RawData from data.Map. This function is
// utility method for other plug-in.
func ConvertMapToRawData(dm data.Map) (RawData, error) {
	// TODO format error, RawData is supposed for cv::Mat structure.
	var width int64
	if w, err := dm.Get(widthPath); err != nil {
		return RawData{}, err
	} else if width, err = data.ToInt(w); err != nil {
		return RawData{}, err
	}

	var height int64
	if h, err := dm.Get(heightPath); err != nil {
		return RawData{}, err
	} else if height, err = data.ToInt(h); err != nil {
		return RawData{}, err
	}

	var img []byte
	if b, err := dm.Get(imagePath); err != nil {
		return RawData{}, err
	} else if img, err = data.ToBlob(b); err != nil {
		return RawData{}, err
	}

	return RawData{
		Width:  int(width),
		Height: int(height),
		Data:   img,
	}, nil
}

// ConvertToDataMap returns data.map. This function is utility method for
// other plug-in.
func (r *RawData) ConvertToDataMap() data.Map {
	// TODO format error
	return data.Map{
		"format": data.String("cvmat"),
		"width":  data.Int(r.Width),
		"height": data.Int(r.Height),
		"image":  data.Blob(r.Data),
	}
}

// ToJpegData convert JPGE format image bytes.
func (r *RawData) ToJpegData(quality int) ([]byte, error) {
	// BGR to RGB
	rgba := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	for i, j := 0, 0; i < len(rgba.Pix); i, j = i+4, j+3 {
		rgba.Pix[i+0] = r.Data[j+2]
		rgba.Pix[i+1] = r.Data[j+1]
		rgba.Pix[i+2] = r.Data[j+0]
		rgba.Pix[i+3] = 0xFF
	}
	w := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
	return w.Bytes(), err
}

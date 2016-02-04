package plugin

import (
	"pfi/sensorbee/opencv"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/udf"
)

// initialize scouter components. this init method will be called by
// SensorBee customized main.go.
//
//  import(
//      _ "pfi/sensorbee/opencv/plugin"
//  )
func init() {
	// capture
	bql.MustRegisterGlobalSourceCreator("opencv_capture_from_uri",
		&opencv.FromURICreator{})
	bql.MustRegisterGlobalSourceCreator("opencv_capture_from_device",
		&opencv.FromDeviceCreator{})

	// cascade classifier
	udf.MustRegisterGlobalUDSCreator("opencv_cascade_classifier",
		udf.UDSCreatorFunc(opencv.NewCascadeClassifier))
	udf.MustRegisterGlobalUDF("opencv_detect_multi_scale",
		udf.MustConvertGeneric(opencv.DetectMultiScale))
}

package plugin

import (
	"pfi/sensorbee/opencv"
	"pfi/sensorbee/sensorbee/bql"
)

// initialize scouter components. this init method will be called by
// SensorBee customized main.go.
//
//  import(
//      _ "pfi/sensorbee/opencv/plugin"
//  )
//
func init() {
	bql.MustRegisterGlobalSourceCreator("opencv_capture_from_uri",
		&opencv.FromURICreator{RawMode: false})
	bql.MustRegisterGlobalSourceCreator("opencv_capture_from_device",
		&opencv.FromDeviceCreator{})

	bql.MustRegisterGlobalSourceCreator("opencv_capture_from_uri_raw",
		&opencv.FromURICreator{RawMode: true})
}

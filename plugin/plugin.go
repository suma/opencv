package plugin

import (
	"fmt"
	"pfi/sensorbee/scouter/capture"
	"pfi/sensorbee/sensorbee/bql"
)

// initialize scouter components. this init method will be called by
// SensorBee customized main.go.
//
//  import(
//      _ "pfi/scouter-snippets/plugin"
//  )
//
// Usage:
//  TYPE capture_from_uri
//    source component, generate frame data from URI
//    (e.g. network camera, video file)
//  TYPE capture_from_device
//    source component, generate frame data from device
func init() {
	sources := []PluginSourceCreator{
		&capture.CaptureFromURI{},
		&capture.CaptureFromDevice{},
	}
	for _, source := range sources {
		creator := source.GetSourceCreator()
		if err := bql.RegisterSourceType(source.TypeName(), creator); err != nil {
			fmt.Errorf("capture plugin registration error: %v", err.Error())
		}
	}
}

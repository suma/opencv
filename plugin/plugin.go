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
	// sources
	sources := []SourceCreator{
		&opencv.FromURICreator{},
		&opencv.FromDeviceCreator{},
	}
	for _, source := range sources {
		bql.MustRegisterGlobalSourceCreator(source.TypeName(), source)
	}
}

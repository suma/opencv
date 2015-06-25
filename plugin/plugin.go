package plugin

import (
	"pfi/sensorbee/scouter/capture"
	"pfi/sensorbee/scouter/detector"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/udf"
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
	// sources
	sources := []PluginSourceCreator{
		&capture.CaptureFromURI{},
		&capture.CaptureFromDevice{},
	}
	for _, source := range sources {
		creator := source.GetSourceCreator()
		if err := bql.RegisterSourceType(source.TypeName(), creator); err != nil {
			panic(err)
		}
	}

	// states
	states := []PluginStateCreator{
		&detector.CameraParameterState{},
	}
	for _, state := range states {
		if err := udf.RegisterGlobalUDSCreator(
			state.TypeName(), udf.UDSCreatorFunc(state.NewState)); err != nil {
			panic(err)
		}
		if err := udf.RegisterGlobalUDF(state.TypeName(), udf.UnaryFunc(state.Func)); err != nil {
			panic(err)
		}
	}
}

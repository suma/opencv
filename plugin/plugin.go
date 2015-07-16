package plugin

import (
	"pfi/sensorbee/scouter/capture"
	"pfi/sensorbee/scouter/detector"
	"pfi/sensorbee/scouter/integrator"
	"pfi/sensorbee/scouter/mjpegserv"
	"pfi/sensorbee/scouter/recognizer"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/udf"
)

// initialize scouter components. this init method will be called by
// SensorBee customized main.go.
//
//  import(
//      _ "pfi/sensorbee/scouter/plugin"
//  )
//
// Usage:
// Source
//  TYPE capture_from_uri
//    source component, generate frame data from URI
//    (e.g. network camera, video file)
//  TYPE capture_from_device
//    source component, generate frame data from device
//
// Sink (TODO)
// State (TODO)
// UDF (TODO)
func init() {
	// sources
	sources := []PluginSourceCreator{
		&capture.CaptureFromURI{},
		&capture.CaptureFromDevice{},
	}
	for _, source := range sources {
		bql.MustRegisterGlobalSourceCreator(source.TypeName(), source)
	}

	// sinks
	mjpegSink := &mjpegserv.MJPEGServ{}
	bql.MustRegisterGlobalSinkCreator("mjpeg_server", mjpegSink)

	// states
	states := []PluginStateCreator{
		&detector.CameraParamState{},
		&detector.ACFDetectionParamState{},
		&detector.MMDetectionParamState{},
		&recog.ImageTaggerCaffeParamState{},
		&integrator.TrackerParamState{},
		&integrator.InstanceManagerParamState{},
	}
	for _, state := range states {
		if err := udf.RegisterGlobalUDSCreator(
			state.TypeName(), udf.UDSCreatorFunc(state.NewState)); err != nil {
			panic(err)
		}
	}

	// UDFs
	udfuncs := []PluginUDFCreator{
		&detector.FrameApplierFuncCreator{},
		&detector.FilterByMaskFuncCreator{},
		&detector.EstimateHeightFuncCreator{},
		&detector.DrawDetectionResultFuncCreator{},
		&detector.FilterByMaskMMFuncCreator{},
		&detector.EstimateHeightMMFuncCreator{},
		&recog.RegionCropFuncCreator{},
	}
	for _, f := range udfuncs {
		udf.MustRegisterGlobalUDF(f.TypeName(), udf.MustConvertGeneric(f.CreateFunction()))
	}

	// UDSFs
	udsfuncs := []PluginUDSFCreator{
		&detector.DetectRegionStreamFuncCreator{},
		&detector.MMDetectRegionStreamFuncCreator{},
		&recog.PredictTagsBatchStreamFuncCreator{},
		&integrator.MovingMatcherStreamFuncCreator{},
		&integrator.FramesTrackerStreamFuncCreator{},
	}
	for _, f := range udsfuncs {
		udf.MustRegisterGlobalUDSFCreator(f.TypeName(), udf.MustConvertToUDSFCreator(f.CreateStreamFunction()))
	}
}

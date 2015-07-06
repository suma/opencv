package plugin

import (
	"pfi/sensorbee/scouter/capture"
	"pfi/sensorbee/scouter/detector"
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
		if err := bql.RegisterGlobalSourceCreator(source.TypeName(), source); err != nil {
			panic(err)
		}
	}

	// sinks
	mjpegSink := &mjpegserv.MJPEGServ{}
	if err := bql.RegisterGlobalSinkCreator("mjpeg_server", mjpegSink); err != nil {
		panic(err)
	}

	// states
	states := []PluginStateCreator{
		&detector.CameraParamState{},
		&detector.ACFDetectionParamState{},
		&recog.ImageTaggerCaffeParamState{},
	}
	for _, state := range states {
		if err := udf.RegisterGlobalUDSCreator(
			state.TypeName(), udf.UDSCreatorFunc(state.NewState)); err != nil {
			panic(err)
		}
	}

	// UDFs
	if err := udf.RegisterGlobalUDF("frame_applier",
		udf.MustConvertGeneric(detector.FrameApplierFunc)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("acf_detector",
		udf.MustConvertGeneric(detector.ACFDetectFunc)); err != nil {
		panic(err)
	}
	udf.RegisterGlobalUDSFCreator("acf_detector_stream", udf.MustConvertToUDSFCreator(detector.CreateACFDetectUDSF))
	if err := udf.RegisterGlobalUDF("filter_by_mask",
		udf.MustConvertGeneric(detector.FilterByMaskFunc)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("estimate_height",
		udf.MustConvertGeneric(detector.EstimateHeightFunc)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("draw_detection_result",
		udf.MustConvertGeneric(detector.DrawDetectionResultFunc)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("recognize_caffe",
		udf.MustConvertGeneric(recog.RecognizeFunc)); err != nil {
		panic(err)
	}

}

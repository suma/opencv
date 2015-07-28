package plugin

import (
	"pfi/sensorbee/scouter/capture"
	"pfi/sensorbee/scouter/detector"
	"pfi/sensorbee/scouter/integrator"
	"pfi/sensorbee/scouter/mjpegserv"
	"pfi/sensorbee/scouter/recognizer"
	"pfi/sensorbee/scouter/utils"
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
func init() {
	// sources
	sources := []PluginSourceCreator{
		&capture.CaptureFromURICreator{},
		&capture.CaptureFromDeviceCreator{},
	}
	for _, source := range sources {
		bql.MustRegisterGlobalSourceCreator(source.TypeName(), source)
	}

	// sinks
	sinks := []PluginSinkCreator{
		&mjpegserv.MJPEGServCreator{},
		&mjpegserv.DebugJPEGWriterCreator{},
	}
	for _, sink := range sinks {
		bql.MustRegisterGlobalSinkCreator(sink.TypeName(), sink)
	}

	// states
	states := []PluginStateCreator{
		&detector.CameraParamState{},
		&detector.ACFDetectionParamState{},
		&detector.MMDetectionParamState{},
		&recog.ImageTaggerCaffeParamState{},
		&integrator.TrackerParamState{},
		&integrator.InstanceManagerParamState{},
	}
	for _, s := range states {
		udf.MustRegisterGlobalUDSCreator(s.TypeName(), udf.UDSCreatorFunc(s.CreateNewState()))
	}

	// UDFs
	udfuncs := []PluginUDFCreator{
		&detector.FrameApplierFuncCreator{},
		&detector.ACFDetectBatchFuncCreator{},
		&detector.FilterByMaskBatchFuncCreator{},
		&detector.EstimateHeightBatchFuncCreator{},
		&detector.FilterByMaskFuncCreator{},
		&detector.EstimateHeightFuncCreator{},
		&detector.DrawDetectionResultFuncCreator{},
		&detector.MMDetectBatchFuncCreator{},
		&detector.FilterByMaskMMBatchFuncCreator{},
		&detector.EstimateHeightMMBatchFuncCreator{},
		&detector.FilterByMaskMMFuncCreator{},
		&detector.EstimateHeightMMFuncCreator{},
		&recog.RegionCropFuncCreator{},
		&recog.PredictTagsFuncCreator{},
		&recog.CroppingAndPredictTagsFuncCreator{},
		&recog.DrawDeteciontResultFuncCreator{},
		&utils.TypeCheckedAggregateFuncCreator{},
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

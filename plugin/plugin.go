package plugin

import (
	"pfi/sensorbee/scouter/capture"
	"pfi/sensorbee/scouter/detector"
	"pfi/sensorbee/scouter/integrator"
	"pfi/sensorbee/scouter/recognizer"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/scouter/writer"
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
	sources := []SourceCreator{
		&capture.FromURICreator{},
		&capture.FromDeviceCreator{},
	}
	for _, source := range sources {
		bql.MustRegisterGlobalSourceCreator(source.TypeName(), source)
	}

	// sinks
	sinks := []SinkCreator{
		&writer.MJPEGServCreator{},
		&writer.JPEGWriterCreator{},
		&writer.VideoWiterCreator{},
	}
	for _, sink := range sinks {
		bql.MustRegisterGlobalSinkCreator(sink.TypeName(), sink)
	}

	// states
	states := []StateCreator{
		&detector.FrameProcessorParamState{},
		&detector.ACFDetectionParamState{},
		&detector.MMDetectionParamState{},
		&recog.ImageTaggerCaffeParamState{},
		&integrator.TrackerParamState{},
		&integrator.InstanceManagerParamState{},
		&integrator.InstancesVisualizerParamState{},
	}
	for _, s := range states {
		udf.MustRegisterGlobalUDSCreator(s.TypeName(),
			udf.UDSCreatorFunc(s.CreateNewState()))
	}

	// UDFs
	udfuncs := []UDFCreator{
		&detector.FrameApplierFuncCreator{},
		&detector.ACFDetectBatchFuncCreator{},
		&detector.FilterByMaskBatchFuncCreator{},
		&detector.EstimateHeightBatchFuncCreator{},
		&detector.PutFeatureBatchUDFCreator{},
		&detector.FilterByMaskFuncCreator{},
		&detector.EstimateHeightFuncCreator{},
		&detector.PutFeatureUDFCreator{},
		&detector.DrawDetectionResultFuncCreator{},
		&detector.MMDetectBatchFuncCreator{},
		&detector.FilterByMaskMMBatchFuncCreator{},
		&detector.EstimateHeightMMBatchFuncCreator{},
		&detector.FilterByMaskMMFuncCreator{},
		&detector.EstimateHeightMMFuncCreator{},
		&recog.CroppingAndPredictTagsFuncCreator{},
		&recog.CroppingAndPredictTagsBatchFuncCreator{},
		&recog.DrawDeteciontResultFuncCreator{},
		&integrator.MultiPlacesMovingMatcherBatchUDFCreator{},
		&integrator.FramesTrackerCacheUDFCreator{},
		&integrator.TrackInstanceStatesUDFCreator{},
		&integrator.InstancesConvertForKanohiJSONUDFCreator{},
		&utils.ObjectCandidateConverterUDFCreator{},
	}
	for _, f := range udfuncs {
		udf.MustRegisterGlobalUDF(f.TypeName(),
			udf.MustConvertGeneric(f.CreateFunction()))
	}

	// UDSFs
	udsfuncs := []UDSFCreator{
		&detector.DetectRegionStreamFuncCreator{},
		&detector.MMDetectRegionStreamFuncCreator{},
		&recog.PredictTagsBatchStreamFuncCreator{},
		&integrator.MultiPlacesMovingMatcherUDSFCreator{},
	}
	for _, f := range udsfuncs {
		udf.MustRegisterGlobalUDSFCreator(f.TypeName(),
			udf.MustConvertToUDSFCreator(f.CreateStreamFunction()))
	}
}

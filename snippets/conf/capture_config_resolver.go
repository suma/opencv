package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
)

const (
	CvCapPropFrameWidth  = 3
	CvCapPropFrameHeight = 4
	CvCapPropFps         = 5
)

type CaptureConfig struct {
	FrameProcessorConfig string
	CameraID             int
	URI                  string
	CaptureFromFile      bool
	FrameSkip            int
	Width                int
	Height               int
	TickInterval         int
}

func GetCaptureSnippetConfig(filePath string) (CaptureConfig, error) {
	conf := CaptureConfig{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	var captureConfig ksconf.CaptureSnippet
	err = json.Unmarshal(file, &captureConfig)
	if err != nil {
		return conf, err
	}

	// get scouter::FrameProcessor::Config
	b, err := json.Marshal(captureConfig.FrameProcessor)
	if err != nil {
		return conf, err
	}
	fpc := string(b)

	return CaptureConfig{
		FrameProcessorConfig: fpc,
		CameraID:             captureConfig.CameraID,
		URI:                  captureConfig.URI,
		CaptureFromFile:      captureConfig.CaptureFromFile,
		FrameSkip:            captureConfig.FrameSkip,
		TickInterval:         captureConfig.TickInterval,
	}, nil
}

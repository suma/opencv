package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
	"pfi/scoutor-snippets/snippets"
	"pfi/scoutor-snippets/snippets/bridge"
)

func GetCaptureSnippetConfig(filePath string) (snippets.CaptureConfig, error) {
	conf := snippets.CaptureConfig{}
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
	fpc := bridge.FrameProcessorConfig_New(string(b))

	return snippets.CaptureConfig{
		FrameProcessorConfig: fpc,
		CameraID:             captureConfig.CameraID,
		URI:                  captureConfig.URI,
		CaptureFromFile:      captureConfig.CaptureFromFile,
		FrameSkip:            captureConfig.FrameSkip,
		TickInterval:         captureConfig.TickInterval,
	}, nil
}

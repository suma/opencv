package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
	"pfi/scouter-snippets/snippets/bridge"
)

type DetectSimpleConfig struct {
	DetectorConfig bridge.DetectorConfig
	PlayerFlag     bool
	JpegQuality    int
}

func GetDetectSimpleSnippetConfig(filePath string) (DetectSimpleConfig, error) {
	conf := DetectSimpleConfig{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	var detectConfig ksconf.DetectSnippet
	err = json.Unmarshal(file, &detectConfig)
	if err != nil {
		return conf, err
	}

	// get scouter::Detector::Config
	b, err := json.Marshal(detectConfig.Detector)
	if err != nil {
		return conf, err
	}
	dc := bridge.DetectorConfig_New(string(b))

	return DetectSimpleConfig{
		DetectorConfig: dc,
		PlayerFlag:     detectConfig.Player != nil,
		JpegQuality:    50,
	}, nil
}

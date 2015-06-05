package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
)

// DetectSimpleConfig is parameter of DetectSimple snippet.
type DetectSimpleConfig struct {
	DetectorConfig string
	PlayerFlag     bool
	JpegQuality    int
}

// GetDetectSimpleSnippetConfig crates configuration data reading external file..
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
	dc := string(b)

	return DetectSimpleConfig{
		DetectorConfig: dc,
		PlayerFlag:     detectConfig.Player != nil,
		JpegQuality:    50,
	}, nil
}

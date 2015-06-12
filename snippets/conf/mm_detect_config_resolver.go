package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
)

// DetectSimpleConfig is parameter of MultiModelDetectSimple snippet.
type MultiModelDetectSimpleConfig struct {
	DetectorConfig string
	PlayerFlag     bool
	JpegQuality    int
}

// GetMultiModelDetectSimpleSnippetConfig crates configuration data reading external file..
func GetMultiModelDetectSimpleSnippetConfig(filePath string) (MultiModelDetectSimpleConfig, error) {
	conf := MultiModelDetectSimpleConfig{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	var detectConfig ksconf.MMDetectSnippet
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

	return MultiModelDetectSimpleConfig{
		DetectorConfig: dc,
		PlayerFlag:     detectConfig.Player != nil,
		JpegQuality:    50,
	}, nil
}

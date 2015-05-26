package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
	"pfi/scouter-snippets/snippets/bridge"
)

type RecognizeCaffeConfig struct {
	ConfigTaggers bridge.RecognizeConfigTaggers
	PlayerFlag    bool
}

func GetRecognizeCaffeSnippetConfig(filePath string) (RecognizeCaffeConfig, error) {
	conf := RecognizeCaffeConfig{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	var recogConfig ksconf.RecognizeCaffe
	err = json.Unmarshal(file, &recogConfig)
	if err != nil {
		return conf, err
	}

	// get std::vector<scouter::ImageTaggerCaffe::Config>
	b, err := json.Marshal(recogConfig.Taggers)
	if err != nil {
		return conf, err
	}
	taggers := bridge.RecognizeConfigTaggers_New(string(b))

	return RecognizeCaffeConfig{
		ConfigTaggers: taggers,
		PlayerFlag:    recogConfig.Player != nil,
	}, nil
}

package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
	"pfi/scoutor-snippets/snippets/bridge"
)

type IntegrateConfig struct {
	IntegrateConfig bridge.IntegratorConfig
	PlayerFlag      bool
}

func GetIntegrateConfig(filePath string) (IntegrateConfig, error) {
	conf := IntegrateConfig{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	var integrateConfig ksconf.IntegrateSnippet
	err = json.Unmarshal(file, &integrateConfig)
	if err != nil {
		return conf, err
	}

	// get scouter::Integrate::Config
	b, err := json.Marshal(integrateConfig.Integrator)
	if err != nil {
		return conf, err
	}
	ic := bridge.IntegratorConfig_New(string(b))

	return IntegrateConfig{
		IntegrateConfig: ic,
		PlayerFlag:      integrateConfig.Player != nil,
	}, nil
}

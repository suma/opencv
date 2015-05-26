package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
)

type IntegrateConfig struct {
	IntegrateConfig string
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
	ic := string(b)

	return IntegrateConfig{
		IntegrateConfig: ic,
		PlayerFlag:      integrateConfig.Player != nil,
	}, nil
}

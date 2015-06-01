package snippets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
)

type DataSenderConfig struct {
	DataSender *ksconf.DataSender
	URI        string
}

type DataSender struct {
	Config DataSenderConfig
}

func (ds *DataSender) SetUp(configPath string) error {
	conf, err := getIntegrateConfig(configPath)
	if err != nil {
		return err
	}

	dataSenderConf := conf.DataSender
	uri := fmt.Sprintf("http://%v:%v%v",
		dataSenderConf.Host, dataSenderConf.Port, dataSenderConf.Path)

	ds.Config = DataSenderConfig{
		DataSender: conf.DataSender,
		URI:        uri,
	}
	return nil
}

// kanochi scouter's data sender information is written in
// integrate.json, and this data sender information get from
// integrate.json.
func getIntegrateConfig(configPath string) (ksconf.IntegrateSnippet, error) {
	conf := ksconf.IntegrateSnippet{}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return conf, err
	}

	var integrateConfig ksconf.IntegrateSnippet
	err = json.Unmarshal(file, &integrateConfig)
	if err != nil {
		return conf, err
	}
	return integrateConfig, nil
}

func (ds *DataSender) Write(ctx *core.Context, t *tuple.Tuple) error {
	is, err := t.Data.Get("instance_states")
	if err != nil {
		return err
	}
	instanceStates, err := is.AsString()
	if err != nil {
		return err
	}

	data := tuple.Map{
		"Send": tuple.Map{
			"time":     tuple.Timestamp(t.Timestamp),
			"instance": tuple.String(instanceStates), // TODO Map is better?
		},
	}
	buf, err := tuple.ToString(data)
	if err != nil {
		return err
	}

	_, err = http.Post(ds.Config.URI, "application/json", strings.NewReader(buf))
	return err
}

func (ds *DataSender) Close(ctx *core.Context) error {
	return nil
}

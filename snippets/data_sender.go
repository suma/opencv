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
	FloorID    int
	DataSender *ksconf.DataSender
}

type DataSender struct {
	Config DataSenderConfig
}

func (ds *DataSender) SetUp(configPath string) error {
	conf, err := getIntegrateConfig(configPath)
	if err != nil {
		return err
	}
	ds.Config = DataSenderConfig{
		FloorID:    conf.FloorID,
		DataSender: conf.DataSender,
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
	conf := ds.Config
	dataSenderConf := conf.DataSender
	uri := fmt.Sprintf("http://%v:%v%v",
		dataSenderConf.Host, dataSenderConf.Port, dataSenderConf.Path)

	data := tuple.Map{
		"Send": tuple.Map{
			"time":     tuple.Timestamp(t.Timestamp),
			"instance": tuple.Map{},
		},
	}
	buf, err := tuple.ToString(data)
	if err != nil {
		return err
	}

	_, err = http.Post(uri, "application/json", strings.NewReader(buf))
	return err
}

func (ds *DataSender) Close(ctx *core.Context) error {
	return nil
}

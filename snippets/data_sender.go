package snippets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
)

// DataSenderConfig is configuration for data sender
type DataSenderConfig struct {
	DataSender *ksconf.DataSender
	playerFlag bool
	uri        string
}

// DataSender is sink component to send result data.
type DataSender struct {
	config DataSenderConfig
}

// SetUp prepares send informations.
func (ds *DataSender) SetUp(configPath string) error {
	conf, err := getIntegrateConfig(configPath)
	if err != nil {
		return err
	}

	dataSenderConf := conf.DataSender
	uri := fmt.Sprintf("http://%v:%v%v",
		dataSenderConf.Host, dataSenderConf.Port, dataSenderConf.Path)

	ds.config = DataSenderConfig{
		DataSender: dataSenderConf,
		playerFlag: conf.Player != nil,
		uri:        uri,
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

// Write instance states given by integrator toward HTTP POST.
func (ds *DataSender) Write(ctx *core.Context, t *tuple.Tuple) error {
	if ds.config.playerFlag {
		go outJpeg(t)
	}
	is, err := t.Data.Get("instance_states")
	if err != nil {
		return nil // usually not set because integrate cache several frames
	}
	instanceStates, err := tuple.AsString(is)
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

	_, err = http.Post(ds.config.uri, "application/json", strings.NewReader(buf))
	return err
}

func outJpeg(t *tuple.Tuple) {
	// detect time to use file name
	ti, err := t.Data.Get("detection_time")
	if err != nil {
		return
	}
	timestamp, _ := tuple.AsTimestamp(ti)
	timeStr := timestamp.Format("15:04:05.999999")

	// detect
	de, err := t.Data.Get("detection_draw_result")
	if err != nil {
		return
	}
	detect, _ := tuple.AsBlob(de)
	ioutil.WriteFile(fmt.Sprintf("detect_%v.jpg", timeStr), detect, os.ModePerm)

	// recognize
	re, err := t.Data.Get("recognize_draw_result")
	if err != nil {
		return
	}
	recog, _ := tuple.AsMap(re)
	for k, v := range recog {
		rec, _ := tuple.AsBlob(v)
		ioutil.WriteFile(fmt.Sprintf("recog[%v]_%v.jpg", k, timeStr),
			rec, os.ModePerm)
	}

	// integrate
	itr, err := t.Data.Get("integrate_result")
	if err != nil {
		return
	}
	integrates, _ := tuple.AsArray(itr)
	for i, v := range integrates {
		integ, _ := tuple.AsBlob(v)
		ioutil.WriteFile(fmt.Sprintf("integrate[%d]_%v.jpg", i, timeStr),
			integ, os.ModePerm)
	}
}

// Close this component.
func (ds *DataSender) Close(ctx *core.Context) error {
	return nil
}

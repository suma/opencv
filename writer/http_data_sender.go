package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/client"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// HTTPDataSenderCreator is a creator of sending data using HTTP client.
type HTTPDataSenderCreator struct{}

var (
	hostPath = data.MustCompilePath("host")
	portPath = data.MustCompilePath("port")
	pathPath = data.MustCompilePath("path")
)

// CreateSink creates sink that data send using HTTP. Data is cast to string.
// In SensorBee, `data.Value` is casted string following JSON spec.
func (c *HTTPDataSenderCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	host := "localhost"
	if h, err := params.Get(hostPath); err == nil {
		if host, err = data.AsString(h); err != nil {
			return nil, err
		}
	}

	var port string
	if po, err := params.Get(portPath); err != nil {
		return nil, fmt.Errorf("not found port number: %v", err.Error())
	} else if port, err = data.ToString(po); err != nil {
		return nil, err
	}

	path := ""
	if p, err := params.Get(pathPath); err == nil {
		if path, err = data.AsString(p); err != nil {
			return nil, err
		}
	}
	url := "http://" + host + ":" + port + path

	return &httpDataSenderSink{
		url: url,
		cli: http.DefaultClient,
	}, nil
}

// TypeName returns name.
func (c *HTTPDataSenderCreator) TypeName() string {
	return "data_sender"
}

type httpDataSenderSink struct {
	url string
	cli *http.Client
}

func (s *httpDataSenderSink) Write(ctx *core.Context, t *core.Tuple) error {
	str := t.Data.String()
	req, err := s.newRequest(str)
	if err != nil {
		return err
	}
	resRaw, err := s.cli.Do(req)
	if err != nil {
		return err
	}
	res := client.Response{
		Raw: resRaw,
	}
	if res.IsError() {
		resErr, err := res.Error()
		if err != nil {
			return err
		}
		return fmt.Errorf("response error: %v", resErr.Message)
	}
	return nil
}

func (s *httpDataSenderSink) newRequest(bodyJSON interface{}) (*http.Request, error) {
	var body io.Reader
	if bodyJSON == nil {
		body = nil
	} else {
		bd, err := json.Marshal(bodyJSON)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bd)
	}

	req, err := http.NewRequest("post", s.url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (s *httpDataSenderSink) Close(ctx *core.Context) error {
	return nil
}

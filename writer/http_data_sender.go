package writer

import (
	"fmt"
	"net/http"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/client"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strings"
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
	if str == "" {
		ctx.Log().Debug("tuple's data is empty")
		return nil
	}
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
			return fmt.Errorf("response error: %v, %v", res.Raw.Status, err)
		}
		return fmt.Errorf("response error: %v", resErr.Message)
	}
	return nil
}

func (s *httpDataSenderSink) newRequest(bodyStr string) (*http.Request, error) {
	body := strings.NewReader(bodyStr)

	req, err := http.NewRequest("POST", s.url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (s *httpDataSenderSink) Close(ctx *core.Context) error {
	return nil
}

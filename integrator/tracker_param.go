package integrator

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type TrackerParamState struct {
	t bridge.Tracker
}

func createTrackerParamState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	p, err := params.Get("file")
	if err != nil {
		return nil, err
	}
	path, err := data.AsString(p)
	if err != nil {
		return nil, err
	}

	// read file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	trackerConfig := string(b)
	s := &TrackerParamState{}
	s.t = bridge.NewTracker(trackerConfig)

	return s, nil
}

func (s *TrackerParamState) CreateNewState() func(*core.Context, data.Map) (core.SharedState, error) {
	return createTrackerParamState
}

func (s *TrackerParamState) TypeName() string {
	return "tracker_parameter"
}

func (s *TrackerParamState) Terminate(ctx *core.Context) error {
	s.t.Delete()
	return nil
}

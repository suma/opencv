package integrator

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// TrackerParamState is a shared state used by tracker UDF/UDSF
type TrackerParamState struct {
	t bridge.Tracker
}

var filePath = data.MustCompilePath("file")

func createTrackerParamState(ctx *core.Context, params data.Map) (core.SharedState,
	error) {
	p, err := params.Get(filePath)
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

// CreateNewState creates a state of Tracking parameters. The parameter is
// collected on JSON, see `scouter::TrackerSP::Config`, which is composition of
// acceptable frame distance, forward frame distance, and so on.
//
// Usage of WITH parameter:
//  "file": tracking parameters file path
func (s *TrackerParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createTrackerParamState
}

func (s *TrackerParamState) TypeName() string {
	return "scouter_tracker_param"
}

func (s *TrackerParamState) Terminate(ctx *core.Context) error {
	s.t.Delete()
	return nil
}

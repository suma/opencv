package integrator

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type InstanceStateConverterUDFCreator struct{}

// CreateFunction returns JSON converter from instance states.
//
// Usage:
//  `convert_instance_states_to_json(states, floorID, timestamp)`
//    states   : instance state array (`[]data.Blob`)
//    floorID  : floor id (`data.Int`)
//    timestamp: timestamp (`data.Timestamp`)
func (c *InstanceStateConverterUDFCreator) CreateFunction() interface{} {
	return convertInstanceStatesToJSON
}

func (c *InstanceStateConverterUDFCreator) TypeName() string {
	return "convert_instance_states_to_json"
}

func convertInstanceStatesToJSON(ctx *core.Context, states data.Array,
	floorID int, timestamp time.Time) (string, error) {

	iss, err := convertToCStates(states)
	if err != nil {
		return "", err
	}

	ts := time.Duration(timestamp.UnixNano()) / time.Millisecond
	json := bridge.ConvertInstanceStatesToJSON(iss, floorID, uint64(ts))

	return json, nil
}

func convertToCStates(states data.Array) ([]bridge.InstanceState, error) {

	iss := []bridge.InstanceState{}
	for _, s := range states {
		b, err := data.AsBlob(s)
		if err != nil {
			return nil, err
		}

		is := bridge.DeserializeInstanceState(b)
		iss = append(iss, is)
	}
	return iss, nil
}

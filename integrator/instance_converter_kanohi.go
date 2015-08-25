package integrator

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// InstancesConvertForKanohiJSONUDFCreator is a creator of converter UDF.
type InstancesConvertForKanohiJSONUDFCreator struct{}

// CreateFunction returns JSON converter from instance states, for kanohi tool.
//
// Usage:
//  `scouter_convert_instances_forkanohijson(states, floorID, timestamp)`
//  [states]
//    * type: []byte
//    * instance states array
//  [floorID]
//    * type: int
//    * the ID of floor to determine the camera.
//  [timestamp]
//    * type: data.Timestamp
//    * captured timestamp, will be converted to [us] (uint64)
//
// Return:
//  The JSON text for kanohi tool.
func (c *InstancesConvertForKanohiJSONUDFCreator) CreateFunction() interface{} {
	return convertInstanceStatesToJSON
}

// TypeName returns type name.
func (c *InstancesConvertForKanohiJSONUDFCreator) TypeName() string {
	return "scouter_convert_instances_forkanohijson"
}

func convertInstanceStatesToJSON(ctx *core.Context, states data.Array,
	floorID int, timestamp data.Timestamp) (string, error) {

	iss, err := convertToCStates(states)
	if err != nil {
		return "", err
	}

	ts, err := data.ToInt(timestamp)
	if err != nil {
		return "", err
	}
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

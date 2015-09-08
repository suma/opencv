package utils

import (
	"github.com/ugorji/go/codec"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// InstanceStateConverterUDFCreator is a creator of `scouter::InstanceState`
// converter to map UDF.
type InstanceStateConverterUDFCreator struct{}

// CreateFunction returns map converter from `scouter::IsntanceState`. Map
// structure is followed by scouter's msgpack packing structure.
func (c *InstanceStateConverterUDFCreator) CreateFunction() interface{} {
	return convertInstanceToMap
}

// TypeName returns type name.
func (c *InstanceStateConverterUDFCreator) TypeName() string {
	return "scouter_convert_instance_to_map"
}

// TODO catch cast error
func convertInstanceToMap(ctx *core.Context, state []byte) (data.Map, error) {
	if state == nil || len(state) == 0 {
		return data.Map{}, nil
	}
	var raw []interface{}
	dec := codec.NewDecoderBytes(state, msgpackHandle)
	err := dec.Decode(&raw)
	if err != nil {
		return nil, err
	}

	return convertState(raw), nil
}

func convertState(raw []interface{}) data.Map {
	tagRaw := raw[1].([]interface{})
	tags := convertTags(tagRaw)

	position := convertPosition(raw[2:5])

	detections := data.Array{}
	obcansRaw := raw[5].([]interface{})
	for _, v := range obcansRaw {
		obcanRaw := v.([]interface{})
		can := convertCandidate(obcanRaw)
		detections = append(detections, can)
	}

	return data.Map{
		"id":         data.Int(raw[0].(uint64)),
		"tags":       tags,
		"position":   position,
		"detections": detections,
	}
}

// InstanceStatesConverterUDFCreator is a creator of `scouter::InstanceState`
// converter to map UDF.
type InstanceStatesConverterUDFCreator struct{}

// CreateFunction returns map converter from `scouter::IsntanceState`. Map
// structure is followed by scouter's msgpack packing structure.
func (c *InstanceStatesConverterUDFCreator) CreateFunction() interface{} {
	return convertInstancesToMap
}

// TypeName returns type name.
func (c *InstanceStatesConverterUDFCreator) TypeName() string {
	return "scouter_convert_instances_to_map"
}

// TODO catch cast error
func convertInstancesToMap(ctx *core.Context, states data.Array) (data.Array, error) {
	convertedArray := data.Array{}
	for _, s := range states {
		b, err := data.ToBlob(s)
		if err != nil {
			return nil, err
		}
		if b == nil || len(b) == 0 {
			continue
		}

		var raw []interface{}
		dec := codec.NewDecoderBytes(b, msgpackHandle)
		err = dec.Decode(&raw)
		if err != nil {
			return nil, err
		}

		convertedArray = append(convertedArray, convertState(raw))
	}
	return convertedArray, nil
}

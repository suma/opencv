package utils

import (
	"github.com/ugorji/go/codec"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

var msgpackHandle = &codec.MsgpackHandle{}

// ObjectCandidateConverterUDFCreator is a creator of `scouter::ObjectCandidate`
// converter to JSON UDF.
type ObjectCandidateConverterUDFCreator struct{}

// CreateFunction returns JSON converter from `scouter::ObjectCandidate`. JSON
// structure is flowed by scouter's msgpack packing structure.
func (c *ObjectCandidateConverterUDFCreator) CreateFunction() interface{} {
	return convertObjectCandidateToJSON
}

// TypeName returns type name.
func (c *ObjectCandidateConverterUDFCreator) TypeName() string {
	return "scouter_convert_regions_to_json"
}

func convertObjectCandidateToJSON(ctx *core.Context, region []byte) (data.Map, error) {
	var raw []interface{}
	dec := codec.NewDecoderBytes(region, msgpackHandle)
	err := dec.Decode(&raw)
	if err != nil {
		return nil, err
	}

	bbox := data.Map{
		"x1": data.Int(raw[0].(uint64)),
		"y1": data.Int(raw[1].(uint64)),
		"x2": data.Int(raw[2].(uint64)),
		"y2": data.Int(raw[3].(uint64)),
	}

	tags := data.Array{}
	tagRaw := raw[6].([]interface{})
	for _, v := range tagRaw {
		r := v.([]interface{})
		tag := data.Map{
			"key":   data.String(string(r[0].([]byte))),
			"value": data.String(string(r[1].([]byte))),
			"score": data.Float(r[2].(float64)),
		}
		tags = append(tags, tag)
	}

	point3f := data.Map{
		"x": data.Float(raw[7].(float64)),
		"y": data.Float(raw[8].(float64)),
		"z": data.Float(raw[9].(float64)),
	}

	featureRaw := raw[11].([]uint8)
	features := data.Array{}
	for _, v := range featureRaw {
		features = append(features, data.Float(float32(v)))
	}

	obcan := data.Map{
		"bbox":       bbox,
		"confidence": data.Float(raw[4].(float64)),
		"tags":       tags,
		"camera_id":  data.Int(raw[5].(int64)),
		"position":   point3f,
		"height":     data.Float(raw[10].(float64)),
		"feature":    features,
	}

	return obcan, nil
}

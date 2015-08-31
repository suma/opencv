package utils

import (
	"github.com/ugorji/go/codec"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

var msgpackHandle = &codec.MsgpackHandle{}

// ObjectCandidateConverterUDFCreator is a creator of `scouter::ObjectCandidate`
// converter to map UDF.
type ObjectCandidateConverterUDFCreator struct{}

// CreateFunction returns map converter from `scouter::ObjectCandidate`. Map
// structure is followed by scouter's msgpack packing structure.
func (c *ObjectCandidateConverterUDFCreator) CreateFunction() interface{} {
	return convertObjectCandidateToMap
}

// TypeName returns type name.
func (c *ObjectCandidateConverterUDFCreator) TypeName() string {
	return "scouter_convert_regions_to_map"
}

// TODO catch cast error
func convertObjectCandidateToMap(ctx *core.Context, regions data.Array) (data.Array,
	error) {
	mapArray := make(data.Array, len(regions))
	for i, region := range regions {
		b, err := data.ToBlob(region)
		if err != nil {
			return nil, err
		}
		var raw []interface{}
		dec := codec.NewDecoderBytes(b, msgpackHandle)
		err = dec.Decode(&raw)
		if err != nil {
			return nil, err
		}

		obcan := convertCandidate(raw)
		mapArray[i] = obcan
	}
	return mapArray, nil
}

func convertCandidate(raw []interface{}) data.Map {
	var x1 int
	if x1u, ok := raw[0].(uint64); ok {
		x1 = int(x1u)
	} else if x1o, ok := raw[0].(int64); ok {
		x1 = int(x1o)
	}
	var y1 int
	if y1u, ok := raw[0].(uint64); ok {
		y1 = int(y1u)
	} else if y1o, ok := raw[0].(int64); ok {
		y1 = int(y1o)
	}
	var x2 int
	if x2u, ok := raw[0].(uint64); ok {
		x2 = int(x2u)
	} else if x2o, ok := raw[0].(int64); ok {
		x2 = int(x2o)
	}
	var y2 int
	if y2u, ok := raw[0].(uint64); ok {
		y2 = int(y2u)
	} else if y2o, ok := raw[0].(int64); ok {
		y2 = int(y2o)
	}
	bbox := data.Map{
		"x1": data.Int(x1),
		"y1": data.Int(y1),
		"x2": data.Int(x2),
		"y2": data.Int(y2),
	}

	tagRaw := raw[6].([]interface{})
	tags := convertTags(tagRaw)

	point3f := convertPosition(raw[7:10])

	obcan := data.Map{
		"bbox":       bbox,
		"confidence": data.Float(raw[4].(float64)),
		"tags":       tags,
		"camera_id":  data.Int(raw[5].(int64)),
		"position":   point3f,
		"height":     data.Float(raw[10].(float64)),
	}

	featureRaw := raw[11].([]uint8)
	if len(featureRaw) != 0 {
		features := data.Array{}
		for _, v := range featureRaw {
			features = append(features, data.Float(float32(v)))
		}
		obcan["feature"] = features
	}
	return obcan
}

func convertTags(raw []interface{}) data.Array {
	tags := data.Array{}
	for _, v := range raw {
		r := v.([]interface{})
		tag := data.Map{
			"key":   data.String(string(r[0].([]byte))),
			"value": data.String(string(r[1].([]byte))),
			"score": data.Float(r[2].(float64)),
		}
		tags = append(tags, tag)
	}
	return tags
}

func convertPosition(raw []interface{}) data.Map {
	point3f := data.Map{
		"x": data.Float(raw[0].(float64)),
		"y": data.Float(raw[1].(float64)),
		"z": data.Float(raw[2].(float64)),
	}
	return point3f
}

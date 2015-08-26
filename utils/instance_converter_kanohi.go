package utils

import (
	"github.com/ugorji/go/codec"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// InstancesConvertForKanohiMapUDFCreator is a creator of converter UDF.
type InstancesConvertForKanohiMapUDFCreator struct{}

// CreateFunction returns JSON converter from instance states, for kanohi tool.
//
// Usage:
//  `scouter_convert_instances_forkanohi(states, floorID, timestamp)`
//  [states]
//    * type: [][]byte
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
func (c *InstancesConvertForKanohiMapUDFCreator) CreateFunction() interface{} {
	return convertInstanceStatesToKanohiMap
}

// TypeName returns type name.
func (c *InstancesConvertForKanohiMapUDFCreator) TypeName() string {
	return "scouter_convert_instances_forkanohi"
}

func convertInstanceStatesToKanohiMap(ctx *core.Context, states data.Array,
	floorID int, timestamp data.Timestamp) (data.Array, error) {

	stateArray := make(data.Array, len(states))
	for i, v := range states {
		b, err := data.AsBlob(v)
		if err != nil {
			return nil, err
		}

		var raw []interface{}
		dec := codec.NewDecoderBytes(b, msgpackHandle)
		err = dec.Decode(&raw)
		if err != nil {
			return nil, err
		}

		is := convertState(raw)
		ismap, err := convertForKanohiStructure(is, floorID)
		if err != nil {
			return nil, err
		}

		ts, err := data.ToInt(timestamp)
		if err != nil {
			return nil, err
		}
		ret := data.Map{
			"time":      data.Int(ts),
			"instances": ismap,
		}

		stateArray[i] = ret
	}

	return stateArray, nil
}

var (
	idPath    = data.MustCompilePath("id")
	xPath     = data.MustCompilePath("position.x")
	yPath     = data.MustCompilePath("position.y")
	tagsPath  = data.MustCompilePath("tags[:]")
	keyPath   = data.MustCompilePath("key")
	valuePath = data.MustCompilePath("value")
)

func convertForKanohiStructure(is data.Map, floorID int) (data.Map, error) {
	id, err := is.Get(idPath)
	if err != nil {
		return nil, err
	}

	x, err := is.Get(xPath)
	if err != nil {
		return nil, err
	}
	y, err := is.Get(yPath)
	if err != nil {
		return nil, err
	}
	loc := data.Map{
		"x":        x,
		"y":        y,
		"floor_id": data.Int(floorID),
	}

	labels := data.Array{}
	if t, err := is.Get(tagsPath); err != nil {
		return nil, err
	} else if tags, err := data.AsArray(t); err != nil {
		return nil, err
	} else {
		for _, tag := range tags {
			var ta data.Map
			if ta, err = data.AsMap(tag); err != nil {
				return nil, err
			}
			var ks string
			if k, err := ta.Get(keyPath); err != nil {
				return nil, err
			} else if ks, err = data.AsString(k); err != nil {
				return nil, err
			}
			var vs string
			if v, err := ta.Get(valuePath); err != nil {
				return nil, err
			} else if vs, err = data.AsString(v); err != nil {
				return nil, err
			}
			labels = append(labels, data.String(ks+"="+vs))
		}
	}
	return data.Map{
		"id":       id,
		"location": loc,
		"labels":   labels,
	}, nil
}

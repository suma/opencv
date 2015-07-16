package utils

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type TypeCheckedAggregateFuncCreator struct{}

func aggregate(ctx *core.Context, items ...data.Value) (data.Array, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("type_checked_aggregate take more than one time")
	}

	a := data.Array{}
	tmpTypeID := items[0].Type()
	for _, i := range items {
		if tmpTypeID != i.Type() {
			return nil, fmt.Errorf("type_checked_aggregate must be all arguments are same type")
		}
		a = append(a, i)
	}
	return a, nil
}

func (c *TypeCheckedAggregateFuncCreator) CreateFunction() interface{} {
	return aggregate
}

func (c *TypeCheckedAggregateFuncCreator) TypeName() string {
	return "type_checked_aggregate"
}

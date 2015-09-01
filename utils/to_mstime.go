package utils

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// ToMSTimeUDFCreator is a creator of convert millisecond timestamp converter.
type ToMSTimeUDFCreator struct{}

// CreateFunction returns a converting millisecond timestamp function.
func (c *ToMSTimeUDFCreator) CreateFunction() interface{} {
	return toMSTime
}

// TypeName returns type name.
func (c *ToMSTimeUDFCreator) TypeName() string {
	return "scouter_to_mstime"
}

func toMSTime(ctx *core.Context, ti data.Timestamp) (int, error) {
	us, err := data.ToInt(ti)
	if err != nil {
		return 0, err
	}

	ms := us / 1e3
	return int(ms), nil
}

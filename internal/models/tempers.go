package models

import "math/rand/v2"

const (
	TEMPER_MAX_VALUE = 6 // non-inclusive
)

type Tempers struct {
	Woe    int64
	Frolic int64
	Dread  int64
	Malice int64
}

func NewTempers() *Tempers {
	return &Tempers{
		Woe:    rand.Int64N(TEMPER_MAX_VALUE),
		Frolic: rand.Int64N(TEMPER_MAX_VALUE),
		Dread:  rand.Int64N(TEMPER_MAX_VALUE),
		Malice: rand.Int64N(TEMPER_MAX_VALUE),
	}
}

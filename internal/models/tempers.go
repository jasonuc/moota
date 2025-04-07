package models

import "math/rand/v2"

const (
	temperMaxVal = 6 // non-inclusive
)

type Tempers struct {
	Woe    int64 // -ve effect on Plant.Xp
	Frolic int64 // +ve effect on Plant.Xp
	Dread  int64 // -ve effect on Plant.Health
	Malice int64 // -ve effect on Plant.Health for nearby plants
}

func NewTempers() *Tempers {
	return &Tempers{
		Woe:    rand.Int64N(temperMaxVal),
		Frolic: rand.Int64N(temperMaxVal),
		Dread:  rand.Int64N(temperMaxVal),
		Malice: rand.Int64N(temperMaxVal),
	}
}

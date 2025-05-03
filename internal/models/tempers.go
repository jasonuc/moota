package models

import "math/rand/v2"

const (
	temperMaxVal = 5
)

// Range: [0, 5]
type Tempers struct {
	Woe    float64 `json:"woe"`    // -ve effect on Plant.XP
	Frolic float64 `json:"frolic"` // +ve effect on Plant.XP
	Dread  float64 `json:"dread"`  // -ve effect on Plant.Health
	Malice float64 `json:"malice"` // -ve effect on Plant.Health for nearby plants
}

func NewTempers() *Tempers {
	return &Tempers{
		Woe:    float64(rand.Int64N(temperMaxVal)) + 1.0,
		Frolic: float64(rand.Int64N(temperMaxVal)) + 1.0,
		Dread:  float64(rand.Int64N(temperMaxVal)) + 1.0,
		Malice: float64(rand.Int64N(temperMaxVal)) + 1.0,
	}
}

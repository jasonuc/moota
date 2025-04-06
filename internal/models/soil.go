package models

import (
	"math/rand/v2"
)

type SoilType int

const (
	SoilTypeLoam SoilType = iota
	SoilTypeSandy
	SoilTypeSilt
	SoilTypeClay
)

type Soil struct {
	ID             string
	Type           SoilType
	WaterRetention float64
	Nutrients      float64
	CircleMeta
}

// TODO: Make Soild.radiusM definitive values like, Small, Medium and Large

func NewSoil(soilType SoilType, centre Coordinates, radiusM float64) *Soil {
	return &Soil{
		Type: soilType,
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: radiusM,
		},
		WaterRetention: rand.Float64() * 5,
		Nutrients:      rand.Float64() * 5,
	}
}

package models

import (
	"math/rand/v2"
	"time"
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
	CreatedAt      time.Time
	CircleMeta
}

// TODO: Make Soild.radiusM definitive values like, Small, Medium and Large

func NewSoil(soilType SoilType, centre Coordinates, radiusM float64, createdAt time.Time) *Soil {
	return &Soil{
		Type: soilType,
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: radiusM,
		},
		WaterRetention: rand.Float64() * 5,
		Nutrients:      rand.Float64() * 5,
		CreatedAt:      createdAt,
	}
}

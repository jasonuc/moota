package models

import (
	"math"
	"math/rand/v2"
	"time"
)

type SoilType string

const (
	SoilTypeLoam  SoilType = "loam"
	SoilTypeSandy SoilType = "sandy"
	SoilTypeSilt  SoilType = "silt"
	SoilTypeClay  SoilType = "clay"
)

type SoilMeta struct {
	Type           SoilType
	WaterRetention float64
	Nutrients      float64
}

var (
	DefaultSoilMetaLoam = SoilMeta{
		Type:           SoilTypeLoam,
		WaterRetention: 0.55,
		Nutrients:      0.75,
	}
	DefaultSoilMetaSandy = SoilMeta{
		Type:           SoilTypeSandy,
		WaterRetention: 0.25,
		Nutrients:      0.20,
	}
	DefaultSoilMetaSilt = SoilMeta{
		Type:           SoilTypeSilt,
		WaterRetention: 0.65,
		Nutrients:      0.55,
	}
	DefaultSoilMetaClay = SoilMeta{
		Type:           SoilTypeClay,
		WaterRetention: 0.80,
		Nutrients:      0.65,
	}
)

type Soil struct {
	ID        string
	CreatedAt time.Time
	SoilMeta
	CircleMeta
}

const (
	SoilRadiusMSmall  = 8.92  // For ≈250 sq. meters
	SoilRadiusMMedium = 17.84 // For ≈1,000 sq. meters
	SoilRadiusMLarge  = 30.90 // For ≈3,000 sq. meters
)

func NewSmallSizedSoil(soilMeta SoilMeta, centre Coordinates, createdAt time.Time) *Soil {
	return newSoil(soilMeta, centre, createdAt, SoilRadiusMSmall)
}

func NewMediumSizedSoil(soilMeta SoilMeta, centre Coordinates, createdAt time.Time) *Soil {
	return newSoil(soilMeta, centre, createdAt, SoilRadiusMMedium)
}

func NewLargeSizedSoil(soilMeta SoilMeta, centre Coordinates, createdAt time.Time) *Soil {
	return newSoil(soilMeta, centre, createdAt, SoilRadiusMLarge)
}

func newSoil(soilMeta SoilMeta, centre Coordinates, createdAt time.Time, radiusM float64) *Soil {
	randomOffset := math.Round((rand.Float64()-0.5)*0.2*100) / 100 // ≈±0.1
	soilMeta.Nutrients = math.Max(0.05, math.Min(1.00, soilMeta.Nutrients+randomOffset))

	randomOffset = math.Round((rand.Float64()-0.5)*0.2*100) / 100 // ≈±0.1
	soilMeta.WaterRetention = math.Max(0.05, math.Min(1.00, soilMeta.WaterRetention+randomOffset))

	return &Soil{
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: radiusM,
		},
		CreatedAt: createdAt,
		SoilMeta:  soilMeta,
	}
}

func (s *Soil) ContainsFullCircle(cm CircleMeta) bool {
	d := cm.Centre().DistanceM(s.Centre())
	return d+cm.RadiusM() <= s.RadiusM()
}

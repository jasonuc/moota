package models

import (
	"math"
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

type SoilMeta struct {
	Type SoilType
	// Min: 0.00; Max 5.00
	WaterRetention float64
	Nutrients      float64
}

var (
	SoilMetaLoam = SoilMeta{
		Type:           SoilTypeLoam,
		WaterRetention: 3.50,
		Nutrients:      4.00,
	}
	SoilMetaSandy = SoilMeta{
		Type:           SoilTypeSandy,
		WaterRetention: 1.50,
		Nutrients:      2.00,
	}
	SoilMetaSilt = SoilMeta{
		Type:           SoilTypeSilt,
		WaterRetention: 3.00,
		Nutrients:      3.50,
	}
	SoilMetaClay = SoilMeta{
		Type:           SoilTypeClay,
		WaterRetention: 4.50,
		Nutrients:      2.50,
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
	// ±3 Offset
	soilMeta.Nutrients += math.Round((rand.Float64()-0.5)*0.6*100) / 100
	soilMeta.WaterRetention += math.Round((rand.Float64()-0.5)*0.6*100) / 100

	return &Soil{
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: radiusM,
		},
		CreatedAt: createdAt,
		SoilMeta:  soilMeta,
	}
}

func (s *Soil) ContainsCircle(cm CircleMeta) bool {
	d := cm.Centre().DistanceM(s.Centre())
	return d+cm.RadiusM() <= s.RadiusM()
}

package models

import (
	"errors"
	"math"
	"math/rand/v2"
	"time"
)

var (
	ErrSoilNotFound = errors.New("soil not found")
)

type SoilType string

const (
	SoilTypeLoam  SoilType = "loam"
	SoilTypeSandy SoilType = "sandy"
	SoilTypeSilt  SoilType = "silt"
	SoilTypeClay  SoilType = "clay"
)

type SoilMeta struct {
	Type             SoilType
	WaterRetention   float64
	NutrientRichness float64
}

var (
	DefaultSoilMetaLoam = SoilMeta{
		Type:             SoilTypeLoam,
		WaterRetention:   0.55,
		NutrientRichness: 0.75,
	}
	DefaultSoilMetaSandy = SoilMeta{
		Type:             SoilTypeSandy,
		WaterRetention:   0.25,
		NutrientRichness: 0.20,
	}
	DefaultSoilMetaSilt = SoilMeta{
		Type:             SoilTypeSilt,
		WaterRetention:   0.65,
		NutrientRichness: 0.55,
	}
	DefaultSoilMetaClay = SoilMeta{
		Type:             SoilTypeClay,
		WaterRetention:   0.80,
		NutrientRichness: 0.65,
	}
)

type Soil struct {
	ID        string
	CreatedAt time.Time
	SoilMeta
	CircleMeta
}

const (
	SoilRadiusMZero   = 0.0
	SoilRadiusMSmall  = 8.92  // For ≈250 sq. meters
	SoilRadiusMMedium = 17.84 // For ≈1,000 sq. meters
	SoilRadiusMLarge  = 30.90 // For ≈3,000 sq. meters
)

func RandomSoilRadius(filterRadius float64) float64 {
	soilRadii := []float64{
		SoilRadiusMSmall,
		SoilRadiusMMedium,
		SoilRadiusMLarge,
	}

	if filterRadius > 0 {
		filteredSoilRadii := []float64{}
		for _, radius := range soilRadii {
			if radius <= filterRadius {
				filteredSoilRadii = append(filteredSoilRadii, radius)
			}
		}
		soilRadii = filteredSoilRadii
	}

	if len(soilRadii) == 0 {
		return SoilRadiusMZero
	}
	randInt := rand.IntN(len(soilRadii))
	return soilRadii[randInt]
}

func MapToNewSizedSoilFn(radius float64) func(SoilMeta, Coordinates) *Soil {
	fns := map[float64]func(SoilMeta, Coordinates) *Soil{
		SoilRadiusMSmall:  NewSmallSizedSoil,
		SoilRadiusMMedium: NewMediumSizedSoil,
		SoilRadiusMLarge:  NewLargeSizedSoil,
	}
	return fns[radius]
}

func RandomSoilMeta() SoilMeta {
	soilMetas := []SoilMeta{
		DefaultSoilMetaLoam,
		DefaultSoilMetaSandy,
		DefaultSoilMetaSilt,
		DefaultSoilMetaClay,
	}
	randInt := rand.IntN(len(soilMetas))
	return soilMetas[randInt]
}

func NewSmallSizedSoil(soilMeta SoilMeta, centre Coordinates) *Soil {
	return newSoil(soilMeta, centre, SoilRadiusMSmall)
}

func NewMediumSizedSoil(soilMeta SoilMeta, centre Coordinates) *Soil {
	return newSoil(soilMeta, centre, SoilRadiusMMedium)
}

func NewLargeSizedSoil(soilMeta SoilMeta, centre Coordinates) *Soil {
	return newSoil(soilMeta, centre, SoilRadiusMLarge)
}

func newSoil(soilMeta SoilMeta, centre Coordinates, radiusM float64) *Soil {
	randomOffset := math.Round((rand.Float64()-0.5)*0.2*100) / 100 // ≈±0.1
	soilMeta.NutrientRichness = math.Max(0.05, math.Min(1.00, soilMeta.NutrientRichness+randomOffset))

	randomOffset = math.Round((rand.Float64()-0.5)*0.2*100) / 100 // ≈±0.1
	soilMeta.WaterRetention = math.Max(0.05, math.Min(1.00, soilMeta.WaterRetention+randomOffset))

	return &Soil{
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: radiusM,
		},
		SoilMeta: soilMeta,
	}
}

func (s *Soil) ContainsFullCircle(cm CircleMeta) bool {
	d := cm.Centre().DistanceM(s.Centre())
	return d+cm.RadiusM() <= s.RadiusM()
}

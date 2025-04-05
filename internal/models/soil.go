package models

import (
	"math"
	"math/rand/v2"
)

type SoilType int

const (
	TypeLoamy SoilType = iota
	TypeSandy
	TypeSilt
	TypeClay
)

type Soil struct {
	ID             string
	Type           SoilType
	centre         Coordinates
	radiusM        float64
	WaterRetention float64
	Nutrients      float64
}

func NewSoil(soilType SoilType, centre Coordinates, radiusM float64) *Soil {
	return &Soil{
		Type:           soilType,
		centre:         centre,
		radiusM:        radiusM,
		WaterRetention: rand.Float64() * 5,
		Nutrients:      rand.Float64() * 5,
	}
}

func (s *Soil) Centre() Coordinates {
	return s.centre
}

func (s *Soil) RadiusM() float64 {
	return s.radiusM
}

func (s *Soil) Area() float64 {
	return math.Pi * (s.RadiusM() * s.RadiusM())
}

func (s *Soil) ContainsPoint(p Coordinates) bool {
	distance := s.centre.DistanceM(p)
	return distance <= s.RadiusM()
}

func (s *Soil) OverlapsWith(other SpatialObject) bool {
	d := s.centre.DistanceM(other.Centre())
	return d <= s.RadiusM()+other.RadiusM()
}

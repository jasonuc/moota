package models

import "math"

type SoilType int

const (
	TypeLoamy SoilType = iota
	TypeSandy
	TypeSilt
	TypeClay
)

type Soil struct {
	ID      string
	Type    SoilType
	Centre  Coordinates
	RadiusM float64 // radius in metres
}

func (s *Soil) Area() float64 {
	return math.Pi * (s.RadiusM * s.RadiusM)
}

func (s *Soil) ContainsPoint(p Coordinates) bool {
	distance := s.Centre.Distance(p)
	return distance <= s.RadiusM
}

func (s *Soil) OverlapsWith(other Soil) bool {
	d := s.Centre.Distance(other.Centre)
	// the two circles touching is considered an overlap as well
	return d <= s.RadiusM+other.RadiusM
}

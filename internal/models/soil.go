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
	ID             string
	Type           SoilType
	centre         Coordinates
	radiusM        float64
	WaterRetention float64
	Nutrients      float64
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
	distance := s.centre.Distance(p)
	return distance <= s.RadiusM()
}

func (s *Soil) OverlapsWith(other SpatialObject) bool {
	d := s.centre.Distance(other.Centre())
	return d <= s.RadiusM()+other.RadiusM()
}

package models

type Circle interface {
	Centre() Coordinates
	RadiusM() float64 // radius in metres
	OverlapsWith(Circle) bool
	ContainsPoint(Coordinates) bool
}

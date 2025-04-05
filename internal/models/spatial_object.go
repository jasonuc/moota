package models

type SpatialObject interface {
	Centre() Coordinates
	RadiusM() float64 // radius in metres
	OverlapsWith(SpatialObject) bool
	ContainsPoint(Coordinates) bool
}

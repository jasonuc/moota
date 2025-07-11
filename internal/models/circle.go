package models

type Circle interface {
	Centre() Coordinates
	RadiusM() float64 // radius in metres
	ContainsPoint(Coordinates) bool
	OverlapsWith(Circle) bool
}

type CircleMeta struct {
	C Coordinates `json:"centre"`
	R float64     `json:"radiusM"`
}

func NewCircleMeta(centre Coordinates, radiusM float64) CircleMeta {
	return CircleMeta{
		C: centre,
		R: radiusM,
	}
}

func (c CircleMeta) Centre() Coordinates {
	return c.C
}

func (c CircleMeta) RadiusM() float64 {
	return c.R
}

func (c CircleMeta) ContainsPoint(p Coordinates) bool {
	return c.Centre().DistanceM(p) < c.RadiusM()
}

func (c CircleMeta) OverlapsWith(other Circle) bool {
	d := c.Centre().DistanceM(other.Centre())
	return d <= c.RadiusM()+other.RadiusM()
}

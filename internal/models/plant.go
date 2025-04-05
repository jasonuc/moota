package models

type Plant struct {
	ID            string
	Nickname      string
	BotanicalName string
	centre        Coordinates
	radiusM       float64 // interaction radius
	SoilID        string
	Health        float64
}

func (p *Plant) Centre() Coordinates {
	return p.centre
}

func (p *Plant) RadiusM() float64 {
	return p.radiusM
}

func (p *Plant) OverlapsWith(other SpatialObject) bool {
	d := p.centre.Distance(other.Centre())
	return d <= p.RadiusM()+other.RadiusM()
}

func (p *Plant) ContainsPoint(other Coordinates) bool {
	distance := p.centre.Distance(other)
	return distance <= p.RadiusM()
}

package models

import "time"

type Plant struct {
	ID            string
	Nickname      string
	centre        Coordinates
	radiusM       float64 // interaction radius
	Health        float64
	PlantedAt     time.Time
	LastWateredAt time.Time
	Xp            int64
	Soil          *Soil
	Tempers       Tempers
	SeedMeta
}

func (p *Plant) Centre() Coordinates {
	return p.centre
}

func (p *Plant) RadiusM() float64 {
	return p.radiusM
}

func (p *Plant) OverlapsWith(other Circle) bool {
	d := p.centre.DistanceM(other.Centre())
	return d <= p.RadiusM()+other.RadiusM()
}

func (p *Plant) ContainsPoint(other Coordinates) bool {
	distance := p.centre.DistanceM(other)
	return distance <= p.RadiusM()
}

func (p *Plant) IsInside(other Circle) bool {
	d := p.Centre().DistanceM(other.Centre())
	return d+p.RadiusM() <= other.RadiusM()
}


// Reference:
// - Calculate distance, bearing and more between Latitude/Longitude points: https://www.movable-type.co.uk/scripts/latlong.html
package models

import (
	"math"
)

// source: https://nssdc.gsfc.nasa.gov/planetary/factsheet/earthfact.html
const earthRadiusM = 6.378e+06 // in metres

type Coordinates struct {
	Lat float64
	Lng float64
}

func (p Coordinates) LatRad() float64 {
	return p.Lat * (math.Pi / 180)
}

func (p Coordinates) LngRad() float64 {
	return p.Lng * (math.Pi / 180)
}

func (p Coordinates) DistanceM(p2 Coordinates) float64 {
	dLat := p2.LatRad() - p.LatRad()
	dLng := p2.LngRad() - p.LngRad()

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(p.LatRad())*math.Cos(p2.LatRad())*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in metres
	return earthRadiusM * c
}

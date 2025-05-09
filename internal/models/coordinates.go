package models

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// source: https://nssdc.gsfc.nasa.gov/planetary/factsheet/earthfact.html
const earthRadiusM = 6.378e+06 // in metres

type Coordinates struct {
	Lat float64 `json:"Lat"`
	Lon float64 `json:"Lon"`
}

func (p Coordinates) latRad() float64 {
	return p.Lat * (math.Pi / 180)
}

func (p Coordinates) lonRad() float64 {
	return p.Lon * (math.Pi / 180)
}

// source: https://www.movable-type.co.uk/scripts/latlong.html#distance
func (p Coordinates) DistanceM(p2 Coordinates) float64 {
	dLat := p2.latRad() - p.latRad()
	dLon := p2.lonRad() - p.lonRad()

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(p.latRad())*math.Cos(p2.latRad())*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in metres
	return earthRadiusM * c
}

// converts a PostGIS POINT string to Coordinates struct.
func CoordinatesFromPostGIS(pointText string) (Coordinates, error) {
	pointText = strings.TrimPrefix(pointText, "POINT(")
	pointText = strings.TrimSuffix(pointText, ")")

	coords := strings.Split(pointText, " ")
	if len(coords) != 2 {
		return Coordinates{}, errors.New("invalid point format")
	}

	lon, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return Coordinates{}, err
	}

	lat, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{Lat: lat, Lon: lon}, nil
}

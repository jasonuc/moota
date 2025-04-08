package models

import (
	"math"
	"testing"
)

func TestDistanceM(t *testing.T) {
	tolerance := 1e-6

	t.Run("distance between the Empire State building and Statue of Liberty", func(t *testing.T) {
		empireStateBuildingP := Coordinates{Lat: 40.7484, Lng: -73.9857}
		statueOfLibertyP := Coordinates{Lat: 40.6892, Lng: -74.0445}

		got := empireStateBuildingP.DistanceM(statueOfLibertyP)
		exp := 8248.546260825362

		if math.Abs(got-exp) > tolerance {
			t.Errorf("got %f but expected %f", got, exp)
		}
	})

	t.Run("distance between the same point on the earth using Big Ben's coordinates", func(t *testing.T) {
		p1 := Coordinates{Lat: 51.5007, Lng: -0.1246}
		p2 := Coordinates{Lat: 51.5007, Lng: -0.1246}

		got := p1.DistanceM(p2)
		exp := 0.0

		if math.Abs(got-exp) > tolerance {
			t.Errorf("got %f but expected %f", got, exp)
		}
	})

	t.Run("distance between antipodal points using Christ the Redeemer's coordinates", func(t *testing.T) {
		p := Coordinates{Lat: -22.9519, Lng: -43.2105}
		pAntipode := Coordinates{Lat: 22.9519, Lng: 136.7895}

		got := p.DistanceM(pAntipode)
		exp := math.Pi * earthRadiusM

		if math.Abs(got-exp) > tolerance {
			t.Errorf("got %f but expected %f", got, exp)
		}
	})
}

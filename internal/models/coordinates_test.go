package models

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistanceM(t *testing.T) {
	t.Run("distance between the Empire State building and Statue of Liberty", func(t *testing.T) {
		empireStateBuildingP := Coordinates{Lat: 40.7484, Lng: -73.9857}
		statueOfLibertyP := Coordinates{Lat: 40.6892, Lng: -74.0445}

		got := empireStateBuildingP.DistanceM(statueOfLibertyP)
		exp := 8248.546260825362
		assert.InDelta(t, got, exp, 1e-5)
	})

	t.Run("distance between the same point on the earth using Big Ben's coordinates", func(t *testing.T) {
		p1 := Coordinates{Lat: 51.5007, Lng: -0.1246}
		p2 := Coordinates{Lat: 51.5007, Lng: -0.1246}

		assert.InDelta(t, p1.DistanceM(p2), 0, 1e-5)
	})

	t.Run("distance between antipodal points using Christ the Redeemer's coordinates", func(t *testing.T) {
		p := Coordinates{Lat: -22.9519, Lng: -43.2105}
		pAntipode := Coordinates{Lat: 22.9519, Lng: 136.7895}

		got := p.DistanceM(pAntipode)
		exp := math.Pi * earthRadiusM
		assert.InDelta(t, got, exp, 1e-5)
	})
}

func TestCoordinatesFromPostGIS(t *testing.T) {
	t.Run("valid point", func(t *testing.T) {
		pointText := "POINT(10.0 20.0)"
		got, err := CoordinatesFromPostGIS(pointText)
		assert.NoError(t, err)
		exp := Coordinates{Lat: 20.0, Lng: 10.0}
		assert.EqualValues(t, got, exp)
	})

	t.Run("invalid point format", func(t *testing.T) {
		pointText := "POINT(10.0)"
		got, err := CoordinatesFromPostGIS(pointText)
		assert.Error(t, err)
		assert.EqualValues(t, got, Coordinates{})
	})
}

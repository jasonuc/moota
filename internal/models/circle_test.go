package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsPoint(t *testing.T) {
	t.Run("return true if point is the centre of circle", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 10.0, Lon: 20.0}, radiusM: 1000.0}
		p := Coordinates{Lat: 10.0, Lon: 20.0}
		assert.True(t, circle.ContainsPoint(p), "expected point to be inside circle")
	})

	t.Run("return true if point is inside the circle", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 10.0, Lon: 20.0}, radiusM: 1000.0}
		p := Coordinates{Lat: 10.003, Lon: 20.003}
		assert.True(t, circle.ContainsPoint(p), "expected point to be inside circle")
	})

	t.Run("return false if circle does not contain the point", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 10.0, Lon: 20.0}, radiusM: 1000.0}
		p := Coordinates{Lat: 10.01, Lon: 20.01}
		assert.False(t, circle.ContainsPoint(p), "expected point to be outside circle")
	})

	t.Run("return false if point is on the circumference of circle", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 37.7749, Lon: -122.4194}, radiusM: 1000.0}
		p := Coordinates{Lat: 37.7839, Lon: -122.4194}
		assert.False(t, circle.ContainsPoint(p), "expected point on circumference to be considered outside circle")
	})
}

func TestOverlapsWith(t *testing.T) {
	t.Run("return true for concentric circles", func(t *testing.T) {
		center := Coordinates{Lat: 10.0, Lon: 20.0}
		circle1 := CircleMeta{C: center, radiusM: 1000.0}
		circle2 := CircleMeta{C: center, radiusM: 2000.0}
		assert.True(t, circle1.OverlapsWith(circle2), "expected concentric circles to be overlapping")
	})

	t.Run("return false for circles that do not overlap", func(t *testing.T) {
		circle1 := CircleMeta{C: Coordinates{Lat: 10.0, Lon: 20.0}, radiusM: 1000.0}
		circle2 := CircleMeta{C: Coordinates{Lat: 10.025, Lon: 20.0}, radiusM: 1000.0}
		assert.False(t, circle1.OverlapsWith(circle2), "expected circles to not be overlapping")
	})

	t.Run("return true for circles that overlap", func(t *testing.T) {
		circle1 := CircleMeta{C: Coordinates{Lat: 10.0, Lon: 20.0}, radiusM: 1000.0}
		circle2 := CircleMeta{C: Coordinates{Lat: 10.012, Lon: 20.0}, radiusM: 1000.0}
		assert.True(t, circle1.OverlapsWith(circle2), "expected circles to be overlapping")
	})

	t.Run("return true for tangential circles", func(t *testing.T) {
		circle1 := CircleMeta{C: Coordinates{Lat: 37.7749, Lon: -122.4194}, radiusM: 1000.0}
		circle2 := CircleMeta{C: Coordinates{Lat: 37.7749, Lon: -122.39125}, radiusM: 1500.0}
		assert.True(t, circle1.OverlapsWith(circle2), "expected tangential circles to be considered overlapping")
	})
}

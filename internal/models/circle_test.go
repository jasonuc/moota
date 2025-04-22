package models

import (
	"testing"
)

func TestContainsPoint(t *testing.T) {
	t.Run("return true if point is the centre of circle", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 10.0, Lng: 20.0}, radiusM: 1000.0}
		p := Coordinates{Lat: 10.0, Lng: 20.0}

		got := circle.ContainsPoint(p)
		exp := true

		if got != exp {
			distance := circle.Centre().DistanceM(p)
			t.Errorf("got %v but expected %v. Distance calculated: %v, Radius: %v",
				got, exp, distance, circle.RadiusM())
		}
	})

	t.Run("return true if point is inside the circle", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 10.0, Lng: 20.0}, radiusM: 1000.0}
		p := Coordinates{Lat: 10.003, Lng: 20.003}

		got := circle.ContainsPoint(p)
		exp := true

		if got != exp {
			distance := circle.Centre().DistanceM(p)
			t.Errorf("got %v but expected %v. Distance calculated: %v, Radius: %v",
				got, exp, distance, circle.RadiusM())
		}
	})

	t.Run("return false if circle does not contain the point", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 10.0, Lng: 20.0}, radiusM: 1000.0}
		p := Coordinates{Lat: 10.01, Lng: 20.01}

		got := circle.ContainsPoint(p)
		exp := false

		if got != exp {
			distance := circle.Centre().DistanceM(p)
			t.Errorf("got %v but expected %v. Distance calculated: %v, Radius: %v",
				got, exp, distance, circle.RadiusM())
		}
	})

	t.Run("return false if point is on the circumference of circle", func(t *testing.T) {
		circle := CircleMeta{C: Coordinates{Lat: 37.7749, Lng: -122.4194}, radiusM: 1000.0}
		p := Coordinates{Lat: 37.7839, Lng: -122.4194}

		got := circle.ContainsPoint(p)
		exp := false

		if got != exp {
			distance := circle.Centre().DistanceM(p)
			t.Errorf("got %v but expected %v. Distance calculated: %v, Radius: %v",
				got, exp, distance, circle.RadiusM())
		}
	})
}

func TestOverlapsWith(t *testing.T) {
	t.Run("return true for concentric circles", func(t *testing.T) {
		center := Coordinates{Lat: 10.0, Lng: 20.0}
		circle1 := CircleMeta{C: center, radiusM: 1000.0}
		circle2 := CircleMeta{C: center, radiusM: 2000.0}

		got := circle1.OverlapsWith(circle2)
		exp := true

		if got != exp {
			distance := circle1.Centre().DistanceM(circle2.Centre())
			t.Errorf("got %v but expected %v. Distance between centers: %v, Sum of radii: %v",
				got, exp, distance, circle1.RadiusM()+circle2.RadiusM())
		}
	})

	t.Run("return false for circles that do not overlap", func(t *testing.T) {
		circle1 := CircleMeta{C: Coordinates{Lat: 10.0, Lng: 20.0}, radiusM: 1000.0}
		circle2 := CircleMeta{C: Coordinates{Lat: 10.025, Lng: 20.0}, radiusM: 1000.0}

		got := circle1.OverlapsWith(circle2)
		exp := false

		if got != exp {
			distance := circle1.Centre().DistanceM(circle2.Centre())
			t.Errorf("got %v but expected %v. Distance between centers: %v, Sum of radii: %v",
				got, exp, distance, circle1.RadiusM()+circle2.RadiusM())
		}
	})

	t.Run("return true for circles that overlap", func(t *testing.T) {
		circle1 := CircleMeta{C: Coordinates{Lat: 10.0, Lng: 20.0}, radiusM: 1000.0}
		circle2 := CircleMeta{C: Coordinates{Lat: 10.012, Lng: 20.0}, radiusM: 1000.0}

		got := circle1.OverlapsWith(circle2)
		exp := true

		if got != exp {
			distance := circle1.Centre().DistanceM(circle2.Centre())
			t.Errorf("got %v but expected %v. Distance between centers: %v, Sum of radii: %v",
				got, exp, distance, circle1.RadiusM()+circle2.RadiusM())
		}
	})

	t.Run("return true for tangential circles", func(t *testing.T) {
		circle1 := CircleMeta{C: Coordinates{Lat: 37.7749, Lng: -122.4194}, radiusM: 1000.0}
		circle2 := CircleMeta{C: Coordinates{Lat: 37.7749, Lng: -122.39125}, radiusM: 1500.0}

		got := circle1.OverlapsWith(circle2)
		exp := true

		if got != exp {
			distance := circle1.Centre().DistanceM(circle2.Centre())
			t.Errorf("got %v but expected %v. Distance between centers: %v, Sum of radii: %v",
				got, exp, distance, circle1.RadiusM()+circle2.RadiusM())
		}
	})
}

package models

import (
	"testing"
)

func TestContainsFullCircle(t *testing.T) {
	t.Run("return true if circle and soil are concentric", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				centre:  Coordinates{Lat: 0, Lng: 0},
				radiusM: 20,
			},
		}

		circle := CircleMeta{
			centre:  Coordinates{Lat: 0, Lng: 0},
			radiusM: 10,
		}

		got := soil.ContainsFullCircle(circle)
		exp := true

		if got != exp {
			t.Errorf("got %v but expected %v", got, exp)
		}
	})

	t.Run("return true if circle is inside soil", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				centre:  Coordinates{Lat: 0, Lng: 0},
				radiusM: 20,
			},
		}

		circle := CircleMeta{
			centre:  Coordinates{Lat: 0.000045, Lng: 0}, // ~5m north at equator
			radiusM: 10,
		}

		got := soil.ContainsFullCircle(circle)
		exp := true

		if got != exp {
			t.Errorf("got %v but expected %v", got, exp)
		}
	})

	t.Run("return false if circle is partially extends outside of soil", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				centre:  Coordinates{Lat: 0, Lng: 0},
				radiusM: 20,
			},
		}

		circle := CircleMeta{
			centre:  Coordinates{Lat: 0.000135, Lng: 0},
			radiusM: 10,
		}

		got := soil.ContainsFullCircle(circle)
		exp := false

		if got != exp {
			t.Errorf("got %v but expected %v", got, exp)
		}
	})

	t.Run("return false if circle is completely outside soil", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				centre:  Coordinates{Lat: 0, Lng: 0},
				radiusM: 20,
			},
		}

		circle := CircleMeta{
			centre:  Coordinates{Lat: 0.00036, Lng: 0},
			radiusM: 10,
		}

		got := soil.ContainsFullCircle(circle)
		exp := false

		if got != exp {
			t.Errorf("got %v but expected %v", got, exp)
		}
	})

	t.Run("return true if circle is tangential to soil but still inside it", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				centre:  Coordinates{Lat: 0.0, Lng: 0.0},
				radiusM: 50,
			},
		}

		circle := CircleMeta{
			centre:  Coordinates{Lat: 0.00026947, Lng: 0.0},
			radiusM: 20,
		}

		got := soil.ContainsFullCircle(circle)
		exp := true

		if got != exp {
			t.Errorf("got %v but expected %v", got, exp)
		}
	})
}

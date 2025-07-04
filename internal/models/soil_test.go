package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsFullCircle(t *testing.T) {
	t.Run("return true if circle and soil are concentric", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				C: Coordinates{Lat: 0, Lon: 0},
				R: 20,
			},
		}

		circle := CircleMeta{
			C: Coordinates{Lat: 0, Lon: 0},
			R: 10,
		}

		assert.True(t, soil.ContainsFullCircle(circle))
	})

	t.Run("return true if circle is inside soil", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				C: Coordinates{Lat: 0, Lon: 0},
				R: 20,
			},
		}

		circle := CircleMeta{
			C: Coordinates{Lat: 0.000045, Lon: 0}, // ~5m north at equator
			R: 10,
		}

		assert.True(t, soil.ContainsFullCircle(circle))
	})

	t.Run("return false if circle is partially extends outside of soil", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				C: Coordinates{Lat: 0, Lon: 0},
				R: 20,
			},
		}

		circle := CircleMeta{
			C: Coordinates{Lat: 0.000135, Lon: 0},
			R: 10,
		}

		assert.False(t, soil.ContainsFullCircle(circle))
	})

	t.Run("return false if circle is completely outside soil", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				C: Coordinates{Lat: 0, Lon: 0},
				R: 20,
			},
		}

		circle := CircleMeta{
			C: Coordinates{Lat: 0.00036, Lon: 0},
			R: 10,
		}

		assert.False(t, soil.ContainsFullCircle(circle))
	})

	t.Run("return true if circle is tangential to soil but still inside it", func(t *testing.T) {
		soil := &Soil{
			CircleMeta: CircleMeta{
				C: Coordinates{Lat: 0.0, Lon: 0.0},
				R: 50,
			},
		}

		circle := CircleMeta{
			C: Coordinates{Lat: 0.00026947, Lon: 0.0},
			R: 20,
		}

		assert.True(t, soil.ContainsFullCircle(circle))
	})
}

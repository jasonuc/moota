package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	simTime := time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC)

	timePlanted := simTime.Add(-24 * time.Hour)
	lastWatered := simTime.Add(-7 * time.Hour)
	lastAction := simTime.Add(-7 * time.Hour)

	plant := &Plant{
		Hp:                100.0,
		Soil:              &Soil{SoilMeta: DefaultSoilMetaLoam},
		LevelMeta:         NewLeveLMeta(1, 0),
		TimePlanted:       timePlanted,
		LastWateredAt:     lastWatered,
		LastActionAt:      lastAction,
		LastRefreshedAt:   nil,
		GracePeriodEndsAt: nil,
	}

	t.Run("water plant after cooldown", func(t *testing.T) {
		now := simTime.Add(4 * time.Hour)

		alive, err := plant.Action(PlantActionWater, now)

		assert.NoError(t, err)
		assert.True(t, alive)
		assert.Equal(t, plant.LastWateredAt, now)
		assert.Equal(t, plant.LastActionAt, now)
		assert.Equal(t, plant.XP, int64(wateringPlantXpGain))
		assert.NotNil(t, plant.GracePeriodEndsAt)
		assert.Equal(t, now.Add(wateringGracePeriod), *plant.GracePeriodEndsAt)
	})

	t.Run("water plant during cooldown", func(t *testing.T) {
		lastWatered := simTime.Add(-2 * time.Hour)
		lastAction := simTime.Add(-2 * time.Hour)
		plant.LastWateredAt = lastWatered
		plant.LastActionAt = lastAction

		now := simTime

		alive, err := plant.Action(PlantActionWater, now)
		assert.ErrorIs(t, err, ErrPlantInCooldown)
		assert.True(t, alive)
		assert.Equal(t, plant.LastWateredAt, lastWatered)
		assert.Equal(t, plant.LastActionAt, lastAction)
	})

	t.Run("water plant exactly at cooldown boundary", func(t *testing.T) {
		lastWatered := simTime
		plant.LastWateredAt = lastWatered

		now := simTime.Add(wateringCooldown)

		alive, err := plant.Action(PlantActionWater, now)
		assert.NoError(t, err)
		assert.True(t, alive)
		assert.Equal(t, plant.LastWateredAt, now)
	})

	t.Run("water plant on new plant", func(t *testing.T) {
		freshPlant := &Plant{
			Hp:                100.0,
			TimePlanted:       simTime,
			LastWateredAt:     time.Time{},
			LastActionAt:      simTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: nil,
		}

		now := simTime.Add(1 * time.Hour)

		alive, err := freshPlant.Action(PlantActionWater, now)
		assert.NoError(t, err)
		assert.True(t, alive)
		assert.Equal(t, freshPlant.LastWateredAt, now)
		assert.NotNil(t, freshPlant.GracePeriodEndsAt)
	})
}

func TestDecayMechanics(t *testing.T) {
	baseTime := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)

	t.Run("plant loses 1 HP every 4 hours without grace period", func(t *testing.T) {
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: nil,
		}

		refreshTime := baseTime.Add(8 * time.Hour)
		plant.Refresh(refreshTime)

		assert.Equal(t, 98.0, plant.Hp)
		assert.False(t, plant.Dead)
	})

	t.Run("plant loses no HP during grace period", func(t *testing.T) {
		gracePeriodEnd := baseTime.Add(4 * time.Hour)
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: &gracePeriodEnd,
		}

		refreshTime := baseTime.Add(3 * time.Hour)
		plant.Refresh(refreshTime)

		assert.Equal(t, 100.0, plant.Hp)
		assert.False(t, plant.Dead)
	})

	t.Run("plant loses HP after grace period expires", func(t *testing.T) {
		waterTime := baseTime.Add(2 * time.Hour)
		gracePeriodEnd := waterTime.Add(4 * time.Hour)
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: &gracePeriodEnd,
		}

		refreshTime := baseTime.Add(8 * time.Hour)
		plant.Refresh(refreshTime)

		assert.Equal(t, 99.0, plant.Hp)
		assert.False(t, plant.Dead)
	})

	t.Run("multiple refresh calls don't compound decay", func(t *testing.T) {
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: nil,
		}

		refreshTime := baseTime.Add(4 * time.Hour)
		plant.Refresh(refreshTime)
		firstHP := plant.Hp

		plant.Refresh(refreshTime.Add(1 * time.Minute))
		secondHP := plant.Hp

		assert.Equal(t, firstHP, secondHP)
	})

	t.Run("plant dies after 16+ days without care", func(t *testing.T) {
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: nil,
		}

		refreshTime := baseTime.Add(17 * 24 * time.Hour)
		plant.Refresh(refreshTime)

		assert.True(t, plant.Dead)
		assert.Equal(t, 0.0, plant.Hp)
		assert.NotNil(t, plant.TimeOfDeath)
	})

	t.Run("refresh skipped if less than 5 minutes since last refresh", func(t *testing.T) {
		lastRefresh := baseTime
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime.Add(-4 * time.Hour),
			LastRefreshedAt:   &lastRefresh,
			GracePeriodEndsAt: nil,
		}

		refreshTime := baseTime.Add(2 * time.Minute)
		plant.Refresh(refreshTime)

		assert.Equal(t, lastRefresh, *plant.LastRefreshedAt)
	})
}

func TestGracePeriodMechanics(t *testing.T) {
	baseTime := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)

	t.Run("watering sets grace period for 4 hours", func(t *testing.T) {
		plant := &Plant{
			Hp:                50.0,
			TimePlanted:       baseTime,
			LastWateredAt:     time.Time{},
			LastActionAt:      baseTime,
			GracePeriodEndsAt: nil,
		}

		waterTime := baseTime.Add(1 * time.Hour)
		plant.Action(PlantActionWater, waterTime)

		expectedGraceEnd := waterTime.Add(wateringGracePeriod)
		assert.NotNil(t, plant.GracePeriodEndsAt)
		assert.Equal(t, expectedGraceEnd, *plant.GracePeriodEndsAt)
	})

	t.Run("watering during existing grace period resets it", func(t *testing.T) {
		firstGraceEnd := baseTime.Add(4 * time.Hour)
		plant := &Plant{
			Hp:                50.0,
			TimePlanted:       baseTime,
			LastWateredAt:     baseTime,
			LastActionAt:      baseTime,
			GracePeriodEndsAt: &firstGraceEnd,
		}

		waterTime := baseTime.Add(3 * time.Hour)
		plant.Action(PlantActionWater, waterTime)

		expectedNewGraceEnd := waterTime.Add(wateringGracePeriod)
		assert.Equal(t, expectedNewGraceEnd, *plant.GracePeriodEndsAt)
	})

	t.Run("grace period cleared when plant dies", func(t *testing.T) {
		gracePeriodEnd := baseTime.Add(4 * time.Hour)
		plant := &Plant{
			Hp:                1.0,
			GracePeriodEndsAt: &gracePeriodEnd,
		}

		plant.Die(baseTime)

		assert.Nil(t, plant.GracePeriodEndsAt)
		assert.True(t, plant.Dead)
	})
}

func TestHelperMethods(t *testing.T) {
	baseTime := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)

	t.Run("IsInGracePeriod returns correct status", func(t *testing.T) {
		gracePeriodEnd := baseTime.Add(4 * time.Hour)
		plant := &Plant{GracePeriodEndsAt: &gracePeriodEnd}

		assert.True(t, plant.IsInGracePeriod(baseTime.Add(2*time.Hour)))
		assert.False(t, plant.IsInGracePeriod(baseTime.Add(5*time.Hour)))

		plantNoGrace := &Plant{GracePeriodEndsAt: nil}
		assert.False(t, plantNoGrace.IsInGracePeriod(baseTime))
	})

	t.Run("CanBeWatered respects cooldown", func(t *testing.T) {
		plant := &Plant{
			Dead:          false,
			LastWateredAt: baseTime,
		}

		assert.False(t, plant.CanBeWatered(baseTime.Add(2*time.Hour)))
		assert.True(t, plant.CanBeWatered(baseTime.Add(3*time.Hour)))
		assert.True(t, plant.CanBeWatered(baseTime.Add(4*time.Hour)))

		deadPlant := &Plant{Dead: true}
		assert.False(t, deadPlant.CanBeWatered(baseTime))

		newPlant := &Plant{Dead: false, LastWateredAt: time.Time{}}
		assert.True(t, newPlant.CanBeWatered(baseTime))
	})

	t.Run("TimeUntilNextWatering calculates correctly", func(t *testing.T) {
		plant := &Plant{
			Dead:          false,
			LastWateredAt: baseTime,
		}

		checkTime := baseTime.Add(1 * time.Hour)
		remaining := plant.TimeUntilNextWatering(checkTime)
		expected := wateringCooldown - 1*time.Hour
		assert.Equal(t, expected, remaining)

		checkTime = baseTime.Add(3 * time.Hour)
		remaining = plant.TimeUntilNextWatering(checkTime)
		assert.Equal(t, time.Duration(0), remaining)
	})

	t.Run("TimeUntilGracePeriodEnds calculates correctly", func(t *testing.T) {
		gracePeriodEnd := baseTime.Add(4 * time.Hour)
		plant := &Plant{GracePeriodEndsAt: &gracePeriodEnd}

		checkTime := baseTime.Add(2 * time.Hour)
		remaining := plant.TimeUntilGracePeriodEnds(checkTime)
		expected := 2 * time.Hour
		assert.Equal(t, expected, remaining)

		checkTime = baseTime.Add(5 * time.Hour)
		remaining = plant.TimeUntilGracePeriodEnds(checkTime)
		assert.Equal(t, time.Duration(0), remaining)
	})
}

func TestChangeHp(t *testing.T) {
	t.Run("increase hp within limits", func(t *testing.T) {
		plant := &Plant{Hp: 50.0}

		plant.changeHp(20.0)

		assert.Equal(t, 70.0, plant.Hp)
		assert.False(t, plant.Dead)
	})

	t.Run("decrease hp within limits", func(t *testing.T) {
		plant := &Plant{Hp: 50.0}

		plant.changeHp(-30.0)

		assert.Equal(t, 20.0, plant.Hp)
		assert.False(t, plant.Dead)
		assert.Nil(t, plant.TimeOfDeath)
	})

	t.Run("decrease plant hp to exactly 0", func(t *testing.T) {
		plant := &Plant{Hp: 10.0}

		plant.changeHp(-10.0)

		assert.Equal(t, 0.0, plant.Hp)
		assert.True(t, plant.Dead)
		assert.NotNil(t, plant.TimeOfDeath)
	})

	t.Run("decrease plant hp beyond 0", func(t *testing.T) {
		plant := &Plant{Hp: 10.0}

		plant.changeHp(-15.0)

		assert.Equal(t, 0.0, plant.Hp)
		assert.True(t, plant.Dead)
		assert.NotNil(t, plant.TimeOfDeath)
	})

	t.Run("increase plant hp above 100", func(t *testing.T) {
		plant := &Plant{Hp: 85.0}

		plant.changeHp(20.0)

		assert.Equal(t, 100.0, plant.Hp)
		assert.False(t, plant.Dead)
	})
}

func TestPerfectCareScenario(t *testing.T) {
	baseTime := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)

	t.Run("plant with perfect care never loses HP", func(t *testing.T) {
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastWateredAt:     time.Time{},
			LastActionAt:      baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: nil,
		}

		currentTime := baseTime
		for i := 0; i < 10; i++ {
			currentTime = currentTime.Add(3 * time.Hour)
			plant.Action(PlantActionWater, currentTime)
			plant.Refresh(currentTime.Add(30 * time.Minute))
		}

		assert.False(t, plant.Dead)
		assert.Equal(t, 100.0, plant.Hp)
	})
}

func TestCasualCareScenario(t *testing.T) {
	baseTime := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)

	t.Run("plant with casual care survives with some HP loss", func(t *testing.T) {
		plant := &Plant{
			Hp:                100.0,
			TimePlanted:       baseTime,
			LastWateredAt:     time.Time{},
			LastActionAt:      baseTime,
			LastRefreshedAt:   nil,
			GracePeriodEndsAt: nil,
		}

		currentTime := baseTime
		for i := 0; i < 5; i++ {
			currentTime = currentTime.Add(8 * time.Hour)
			plant.Action(PlantActionWater, currentTime)
			plant.Refresh(currentTime.Add(30 * time.Minute))
		}

		assert.False(t, plant.Dead)
		assert.Greater(t, plant.Hp, 90.0)
	})
}

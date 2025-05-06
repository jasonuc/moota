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
		Hp:              100.0,
		Soil:            &Soil{SoilMeta: DefaultSoilMetaLoam},
		LevelMeta:       NewLeveLMeta(1, 0),
		TimePlanted:     timePlanted,
		LastWateredTime: lastWatered,
		LastActionTime:  lastAction,
	}
	t.Run("water plant after cooldown", func(t *testing.T) {
		now := simTime.Add(2 * time.Hour)

		alive, err := plant.Action(PlantActionWater, now)

		assert.NoError(t, err)
		assert.True(t, alive, "expected plant to be alive")
		assert.EqualValues(t, plant.LastWateredTime, now)
		assert.EqualValues(t, plant.LastActionTime, now)
		assert.Equal(t, plant.XP, int64(wateringPlantXpGain))
	})

	t.Run("Water plant during cooldown", func(t *testing.T) {
		lastWatered := simTime.Add(-1 * time.Hour)
		lastAction := simTime.Add(-1 * time.Hour)
		plant.LastWateredTime = lastWatered
		plant.LastActionTime = lastAction

		now := simTime

		alive, err := plant.Action(PlantActionWater, now)
		assert.ErrorIs(t, err, ErrPlantInCooldown)
		assert.True(t, alive, "expected plant to be alive")
		assert.EqualValues(t, plant.LastWateredTime, lastWatered)
		assert.EqualValues(t, plant.LastActionTime, lastAction)
	})
}

func TestPreActionHook(t *testing.T) {
	t.Run("hp decreases when plant has not been interacted with for over 2 days", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(2 * 24 * time.Hour)
		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: currentTime, LastWateredTime: currentTime}

		plant.preActionHook(futureTime)

		expHp := 36.80
		assert.False(t, plant.Dead, "expected plant to be alive")
		assert.InDelta(t, plant.Hp, expHp, 0.01)
	})

	t.Run("hp decreases when plant has not been watered for over 7 hours", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(13 * time.Hour)
		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: futureTime, LastWateredTime: currentTime}

		plant.preActionHook(futureTime)

		expHp := 49.7
		assert.False(t, plant.Dead, "expected plant to be alive")
		assert.InDelta(t, plant.Hp, expHp, 0.01)
	})

	t.Run("hp decreases down to zero if neglected for a long time", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(7 * 24 * time.Hour)
		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: currentTime, LastWateredTime: currentTime}

		plant.preActionHook(futureTime)

		assert.True(t, plant.Dead)
		assert.Zero(t, plant.Hp)
	})
}

func TestChangeHp(t *testing.T) {
	t.Run("increase hp within limits", func(t *testing.T) {
		plant := &Plant{Hp: 50.0}

		plant.changeHp(20.0)
		expHp := 70.0

		assert.Equal(t, plant.Hp, expHp)
		assert.False(t, plant.Dead)
	})

	t.Run("decrease hp within limits", func(t *testing.T) {
		plant := &Plant{Hp: 50.0}

		plant.changeHp(-30.0)
		expHp := 20.0

		assert.Equal(t, plant.Hp, expHp)
		assert.False(t, plant.Dead)
		assert.Nil(t, plant.TimeOfDeath)
	})

	t.Run("decrease plant hp beyond 0", func(t *testing.T) {
		plant := &Plant{Hp: 10.0}

		plant.changeHp(-15.0)
		expHp := 0.0

		assert.Equal(t, plant.Hp, expHp)
		assert.True(t, plant.Dead)
		assert.NotNil(t, plant.TimeOfDeath)
	})

	t.Run("increase plant hp above 100", func(t *testing.T) {
		plant := &Plant{Hp: 85.0}

		plant.changeHp(20)
		expHp := 100.0

		assert.Equal(t, plant.Hp, expHp)
		assert.False(t, plant.Dead, "expected plant to be dead")
	})
}

func TestRefresh(t *testing.T) {
	t.Run("refresh plant for plant that would die after refresh", func(t *testing.T) {
		timeForRefresh := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
		lastActionTime := time.Date(2025, 4, 24, 12, 0, 0, 0, time.UTC) // 7 days earlier

		plant := &Plant{
			LevelMeta:       LevelMeta{Level: 1, XP: 0},
			Hp:              50.0,
			LastActionTime:  lastActionTime,
			LastWateredTime: lastActionTime,
		}

		plant.Refresh(timeForRefresh)

		assert.True(t, plant.Dead, "expected plant to be dead")
		assert.Zero(t, plant.Hp)
		assert.NotNil(t, plant.TimeOfDeath)
	})

	t.Run("refresh plant for plant that would not die after refresh", func(t *testing.T) {
		baseTime := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
		timeForRefresh := baseTime.Add(36 * time.Hour)
		lastActionTime := baseTime

		plant := &Plant{
			LevelMeta:       LevelMeta{Level: 1, XP: 0},
			Hp:              50.0,
			LastActionTime:  lastActionTime,
			LastWateredTime: lastActionTime,
		}

		plant.Refresh(timeForRefresh)

		assert.False(t, plant.Dead, "expected plant to be alive")
		assert.InDelta(t, 41.6, plant.Hp, 0.01, "expected approximately 41.6 hp")
		assert.Nil(t, plant.TimeOfDeath)
	})
}

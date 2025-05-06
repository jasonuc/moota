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

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !alive {
			t.Errorf("expected plant to be alive")
		}

		if plant.LastWateredTime != now {
			t.Errorf("expected LastWateredTime to be updated to %v, got: %v", now, plant.LastWateredTime)
		}

		if plant.LastActionTime != now {
			t.Errorf("expected LastActionTime to be updated to %v, got: %v", now, plant.LastActionTime)
		}

		if plant.XP != wateringPlantXpGain {
			t.Errorf("expected xp %v, got: %v", wateringPlantXpGain, plant.XP)
		}
	})

	t.Run("Water plant during cooldown", func(t *testing.T) {
		lastWatered := simTime.Add(-1 * time.Hour)
		lastAction := simTime.Add(-1 * time.Hour)
		plant.LastWateredTime = lastWatered
		plant.LastActionTime = lastAction

		now := simTime

		alive, err := plant.Action(PlantActionWater, now)

		if err != ErrPlantInCooldown {
			t.Fatalf("expected ErrInCooldown, got: %v", err)
		}

		if !alive {
			t.Errorf("expected plant to still be alive")
		}

		if plant.LastWateredTime.Equal(now) {
			t.Errorf("LastWateredTime should not have been updated")
		}

		if plant.LastActionTime.Equal(now) {
			t.Errorf("LastWateredTime should not have been updated")
		}
	})
}

func TestPreActionHook(t *testing.T) {
	t.Run("hp decreases when plant has not been interacted with for over 2 days", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(2 * 24 * time.Hour)
		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: currentTime, LastWateredTime: currentTime}

		plant.preActionHook(futureTime)

		expHp := 36.80

		if !plant.Alive() {
			t.Errorf("expected plant to be alive")
		}

		if !assert.InDelta(t, expHp, plant.Hp, 0.01) {
			t.Errorf("expected %f hp but got %f", expHp, plant.Hp)
		}
	})

	t.Run("hp decreases when plant has not been watered for over 7 hours", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(13 * time.Hour)
		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: futureTime, LastWateredTime: currentTime}

		plant.preActionHook(futureTime)

		expHp := 49.7

		if !plant.Alive() {
			t.Errorf("expected plant to be alive")
		}

		if plant.Hp != expHp {
			t.Errorf("expected %f hp but got %f", expHp, plant.Hp)
		}
	})

	t.Run("hp decreases down to zero if neglected for a long time", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(7 * 24 * time.Hour)
		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: currentTime, LastWateredTime: currentTime}

		plant.preActionHook(futureTime)

		expHp := 0.0

		if plant.Alive() {
			t.Errorf("expected plant to be dead")
		}

		if plant.Hp != expHp {
			t.Errorf("expected %f hp but got %f", expHp, plant.Hp)
		}
	})
}

func TestChangeHp(t *testing.T) {
	t.Run("increase hp within limits", func(t *testing.T) {
		plant := &Plant{Hp: 50.0}

		alive := plant.changeHp(20.0)
		expHp := 70.0

		if plant.Hp != expHp {
			t.Errorf("expected hp %f but got %f", expHp, plant.Hp)
		}

		if !alive {
			t.Errorf("expected plant to be alive")
		}
	})

	t.Run("decrease hp within limits", func(t *testing.T) {
		plant := &Plant{Hp: 50.0}

		alive := plant.changeHp(-30.0)
		expHp := 20.0

		if plant.Hp != expHp {
			t.Errorf("expected hp %f but got %f", expHp, plant.Hp)
		}

		if !alive {
			t.Errorf("expected plant to be alive")
		}

		if plant.TimeOfDeath != nil {
			t.Errorf("expected plant to not have a time of death")
		}
	})

	t.Run("decrease plant hp beyond 0", func(t *testing.T) {
		plant := &Plant{Hp: 10.0}

		alive := plant.changeHp(-15.0)
		expHp := 0.0

		if plant.Hp != expHp {
			t.Errorf("expected hp %f but got %f", expHp, plant.Hp)
		}

		if alive {
			t.Errorf("expected plant to be dead")
		}

		if plant.TimeOfDeath == nil {
			t.Errorf("expected plant to have a time of death")
		}
	})

	t.Run("increase plant hp above 100", func(t *testing.T) {
		plant := &Plant{Hp: 85.0}

		alive := plant.changeHp(20)
		expHp := 100.0

		if plant.Hp != expHp {
			t.Errorf("expected hp %f but got %f", plant.Hp, expHp)
		}

		if !alive {
			t.Errorf("expected plant ot be alive")
		}
	})
}

func TestRefresh(t *testing.T) {
	t.Run("refresh plant for plant that would die after refresh", func(t *testing.T) {
		currentTime := time.Now()
		timeForRefresh := currentTime.Add(7 * 24 * time.Hour)

		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: time.Now().Add(-7 * 24 * time.Hour), LastWateredTime: time.Now().Add(-7 * 24 * time.Hour)}

		plant.Refresh(timeForRefresh)

		expAlive := false
		expHp := 0.0
		expTimeOfDeath := timeForRefresh

		if plant.Alive() != expAlive {
			t.Errorf("expected plant to be dead")
		}

		if plant.Hp != expHp {
			t.Errorf("expected %f hp but got %f", expHp, plant.Hp)
		}

		if plant.TimeOfDeath.Equal(expTimeOfDeath) {
			t.Errorf("expected plant to have a time of death")
		}
	})

	t.Run("refresh plant for plant that would not die after refresh", func(t *testing.T) {
		currentTime := time.Now()
		timeForRefresh := currentTime.Add(36 * time.Hour)

		plant := &Plant{LevelMeta: LevelMeta{Level: 1, XP: 0}, Hp: 50.0, LastActionTime: currentTime, LastWateredTime: currentTime}

		plant.Refresh(timeForRefresh)

		expAlive := true
		expHP := 41.6

		if plant.Alive() != expAlive {
			t.Errorf("expected plant to be alive")
		}

		if !assert.InDelta(t, expHP, plant.Hp, 0.01) {
			t.Errorf("expected %f hp but got %f", expHP, plant.Hp)
		}

		if plant.TimeOfDeath != nil {
			t.Errorf("expected plant to not have a time of death")
		}
	})
}

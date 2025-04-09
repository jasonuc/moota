package models

import (
	"testing"
	"time"
)

func TestAction(t *testing.T) {
	simTime := time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC)

	plant := &Plant{
		Hp:             100.0,
		Level:          1,
		Soil:           &Soil{SoilMeta: SoilMetaLoam},
		PlantedAt:      simTime.Add(-24 * time.Hour),
		LastWateredAt:  simTime.Add(-7 * time.Hour),
		LastActionTime: simTime.Add(-7 * time.Hour),
	}

	t.Run("Water plant after cooldown", func(t *testing.T) {
		now := simTime.Add(2 * time.Hour)

		alive, err := plant.Action(PlantActionWater, now)

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !alive {
			t.Errorf("expected plant to be alive")
		}

		if plant.LastWateredAt != now {
			t.Errorf("expected LastWateredAt to be updated to %v, got: %v", now, plant.LastWateredAt)
		}

		if plant.LastActionTime != now {
			t.Errorf("expected LastActionTime to be updated to %v, got: %v", now, plant.LastActionTime)
		}

		if plant.Xp != wateringPlantXpGain {
			t.Errorf("expected xp %v, got: %v", wateringPlantXpGain, plant.Xp)
		}
	})

	t.Run("Water plant during cooldown", func(t *testing.T) {
		plant.LastWateredAt = simTime.Add(-1 * time.Hour)
		plant.LastActionTime = simTime.Add(-1 * time.Hour)

		now := simTime

		alive, err := plant.Action(PlantActionWater, now)

		if err != ErrPlantInCooldown {
			t.Fatalf("expected ErrInCooldown, got: %v", err)
		}

		if !alive {
			t.Errorf("expected plant to still be alive")
		}

		if plant.LastWateredAt.Equal(now) {
			t.Errorf("LastWateredAt should not have been updated")
		}

		if plant.LastActionTime.Equal(now) {
			t.Errorf("LastWateredAt should not have been updated")
		}
	})
}

func TestPreActionHook(t *testing.T) {
	t.Run("hp decreases when plant has not been interacted with for over 12 hours", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(14 * time.Hour)
		plant := &Plant{Level: 1, Xp: 0, Hp: 50.0, LastActionTime: currentTime, LastWateredAt: futureTime}

		plant.preActionHook(futureTime)

		expHp := 40.0

		if !plant.Alive() {
			t.Errorf("expected plant to be alive")
		}

		if plant.Hp != expHp {
			t.Errorf("expected %f hp but got %f", expHp, plant.Hp)
		}
	})

	t.Run("hp decreases when plant has not been watered for over 7 hours", func(t *testing.T) {
		currentTime := time.Now()
		futureTime := currentTime.Add(10 * time.Hour)
		plant := &Plant{Level: 1, Xp: 0, Hp: 50.0, LastActionTime: futureTime, LastWateredAt: currentTime}

		plant.preActionHook(futureTime)

		expHp := 47.0

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
		plant := &Plant{Level: 1, Xp: 0, Hp: 50.0, LastActionTime: currentTime, LastWateredAt: currentTime}

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

func TestAddXpAndLeveUp(t *testing.T) {
	t.Run("add xp but not enough to increase level", func(t *testing.T) {
		plant := &Plant{Level: 1, Xp: 0}
		expPlant := struct {
			Level int64
			Xp    int64
		}{Level: 1, Xp: 25}

		plant.addXp(25)

		if plant.Xp != expPlant.Xp {
			t.Errorf("expected %d xp but got %d", expPlant.Xp, plant.Xp)
		}

		if plant.Level != expPlant.Level {
			t.Errorf("expected level %d but got %d", expPlant.Level, plant.Level)
		}
	})

	t.Run("add enough xp for plant to level up once", func(t *testing.T) {
		plant := &Plant{Level: 1, Xp: 0}
		expPlant := struct {
			Level int64
			Xp    int64
		}{Level: 2, Xp: 25}

		xpToAdd := xpRequiredForLevel(expPlant.Level) + 25
		plant.addXp(xpToAdd)

		if plant.Xp != expPlant.Xp {
			t.Errorf("expected %d xp but got %d", expPlant.Xp, plant.Xp)
		}

		if plant.Level != expPlant.Level {
			t.Errorf("expected level %d but got %d", expPlant.Level, plant.Level)
		}
	})

	t.Run("add enough xp for plant to level up multiple times", func(t *testing.T) {
		plant := &Plant{Level: 1, Xp: 0}
		expPlant := struct {
			Level int64
			Xp    int64
		}{Level: 5, Xp: 10}

		xpToAdd := xpRequiredForLevel(2) + xpRequiredForLevel(3) + xpRequiredForLevel(4) + xpRequiredForLevel(5) + expPlant.Xp
		plant.addXp(xpToAdd)

		if plant.Xp != expPlant.Xp {
			t.Errorf("expected %d xp but got %d", expPlant.Xp, plant.Xp)
		}

		if plant.Level != expPlant.Level {
			t.Errorf("expected level %d but got %d", expPlant.Level, plant.Level)
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

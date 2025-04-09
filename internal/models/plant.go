package models

import (
	"errors"
	"math"
	"time"
)

const (
	plantInteractionRadius = 3.0
	wateringPlantXpGain    = 30 // TODO: Up for debate, meaning `xpRequiredForLevel` functionality might have to be changed in the future
	minWateringInterval    = 6 * time.Hour
)

var (
	ErrPlantInCooldown     = errors.New("error plant in cooldown mode")
	ErrPlantNotFullyInSoil = errors.New("error plant not fully inside soil")
)

type PlantAction int

const (
	PlantActionWater PlantAction = iota
)

type Plant struct {
	ID             string
	Nickname       string
	Hp             float64
	Xp             int64
	Level          int64
	Dead           bool
	Soil           *Soil
	Tempers        *Tempers
	PlantedAt      time.Time
	LastWateredAt  time.Time
	LastActionTime time.Time
	SeedMeta
	CircleMeta
}

func NewPlant(seed *Seed, soil *Soil, centre Coordinates, plantedAt time.Time) (*Plant, error) {
	nickname := "Atura" // TODO: Make a function to generate random whimsical names
	seed.Planted = true
	circleMeta := CircleMeta{radiusM: plantInteractionRadius, centre: centre}
	if !soil.ContainsFullCircle(circleMeta) {
		return nil, ErrPlantNotFullyInSoil
	}

	healthOffset := 0.0
	var xpBonus int64
	if seed.SeedMeta.OptimalSoil == soil.Type {
		healthOffset += 15.0
		xpBonus += 25
	} else {
		if seed.SeedMeta.IsCompatibleWithSoil(soil.Type) {
			healthOffset += 5.0
		} else {
			healthOffset -= 5.0
		}
	}

	return &Plant{
		Nickname:  nickname,
		Hp:        seed.Health + healthOffset,
		Soil:      soil,
		Xp:        xpBonus,
		Level:     1,
		Tempers:   NewTempers(),
		PlantedAt: plantedAt,
		SeedMeta:  seed.SeedMeta,
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: plantInteractionRadius,
		},
	}, nil
}

func (p *Plant) Action(action PlantAction, t time.Time) (bool, error) {
	p.preActionHook(t)

	switch action {
	case PlantActionWater:
		if t.Sub(p.LastWateredAt) > minWateringInterval {
			p.addXp(wateringPlantXpGain)
			p.LastWateredAt = t
		} else {
			return p.Alive(), ErrPlantInCooldown
		}
	}
	p.LastActionTime = t

	return p.Alive(), nil
}

func (p *Plant) preActionHook(t time.Time) {
	// Hp reduction for plant neglect
	hoursSinceLastAction := t.Sub(p.LastActionTime).Hours()
	decreaseMult := math.Floor(hoursSinceLastAction - 12)
	if decreaseMult > 0 {
		p.changeHp(-5 * decreaseMult)
	}

	// Hp reduction for lack of watering
	hoursSinceLastWatering := math.Floor(t.Sub(p.LastWateredAt).Hours())
	decreaseMult = math.Floor(hoursSinceLastWatering - 7)
	if decreaseMult > 0 {
		p.changeHp(-1 * decreaseMult)
	}
}

func (p *Plant) Alive() bool {
	return !p.Dead
}

func xpRequiredForLevel(level int64) int64 {
	return int64(math.Round(75 * (math.Pow(float64(level), 2.0) - float64(level))))
}

func (p *Plant) addXp(xp int64) {
	p.Xp += xp

	// Handles scenario where plant needs to level up multiple times in one go
	for {
		nextLevelXpReq := xpRequiredForLevel(p.Level + 1)
		if p.Xp >= nextLevelXpReq {
			p.Xp -= nextLevelXpReq
			p.levelUp()
		} else {
			break
		}
	}
}

func (p *Plant) levelUp() int64 {
	p.Level += 1
	return p.Level
}

func (p *Plant) changeHp(delta float64) bool {
	p.Hp = math.Max(0, math.Min(100, p.Hp+delta)) // clamp hp between 0 and 100
	if p.Hp == 0 {
		p.Dead = true
	}
	return p.Alive()
}

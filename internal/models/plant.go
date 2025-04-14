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
	ErrPlantInCooldown     = errors.New("plant in cooldown mode")
	ErrPlantNotFullyInSoil = errors.New("plant not fully inside soil")
)

type PlantAction int

const (
	PlantActionWater PlantAction = iota
)

type Plant struct {
	ID             string
	Nickname       string
	Hp             float64
	Dead           bool
	OwnerID        string
	Soil           *Soil
	Tempers        *Tempers
	PlantedAt      time.Time
	LastWateredAt  time.Time
	LastActionTime time.Time
	SeedMeta
	LevelMeta
	CircleMeta
}

func NewPlant(seed *Seed, soil *Soil, centre Coordinates, plantedAt time.Time) (*Plant, error) {
	if seed.Planted {
		return nil, ErrSeedAlreadyPlanted
	}

	nickname := "Atura" // TODO: Make a function to generate random whimsical names
	seed.Planted = true
	circleMeta := CircleMeta{radiusM: plantInteractionRadius, centre: centre}
	if !soil.ContainsFullCircle(circleMeta) {
		return nil, ErrPlantNotFullyInSoil
	}

	healthOffset := 0.0
	var xpBonus int64
	if seed.OptimalSoil == soil.Type {
		healthOffset += 15.0
		xpBonus += 25
	} else {
		if seed.IsCompatibleWithSoil(soil.Type) {
			healthOffset += 5.0
		} else {
			healthOffset -= 5.0
		}
	}

	return &Plant{
		Nickname:  nickname,
		Hp:        seed.Health + healthOffset,
		Soil:      soil,
		OwnerID:   seed.OwnerID,
		Dead:      false,
		LevelMeta: NewLeveLMeta(1, xpBonus),
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
	// TODO: consider converting all time to UTC before carrying out any action
	p.preActionHook(t)

	switch action {
	case PlantActionWater:
		// TODO: move this to it's own function and make use of the Soil.WaterRetention
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

func (p *Plant) changeHp(delta float64) bool {
	p.Hp = math.Max(0, math.Min(100, p.Hp+delta)) // clamp hp between 0 and 100
	if p.Hp == 0 {
		p.Dead = true
	}
	return p.Alive()
}

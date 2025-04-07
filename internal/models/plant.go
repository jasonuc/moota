package models

import (
	"fmt"
	"math"
	"time"
)

const (
	plantInteractionRadius = 2.0 // TODO: This value is up for debate
	plantInitialLevel      = 1
)

type Plant struct {
	ID            string
	Nickname      string
	Health        float64
	Xp            int64
	Level         int64
	Dead          bool
	Soil          *Soil
	Tempers       *Tempers
	PlantedAt     time.Time
	LastWateredAt time.Time
	SeedMeta
	CircleMeta
}

func NewPlant(seed *Seed, soil *Soil, centre Coordinates) (*Plant, error) {
	nickname := "Atura" // TODO: Make a function to generate random whimsical names
	seed.Planted = true
	circleMeta := CircleMeta{radiusM: plantInteractionRadius, centre: centre}
	if !IsInsideSoil(circleMeta, soil) {
		return nil, fmt.Errorf("plant is not completely inside soil")
	}

	healthOffset := 0.0
	var xpBonus int64
	if seed.SeedMeta.OptimalSoil == soil.Type {
		healthOffset += 1.5
		xpBonus += 25
	} else {
		if seed.SeedMeta.IsCompatibleWithSoil(soil.Type) {
			healthOffset += 0.5
		} else {
			healthOffset -= 0.5
		}
	}

	return &Plant{
		Nickname:  nickname,
		Health:    seed.Health + healthOffset,
		Soil:      soil,
		Xp:        xpBonus,
		Level:     plantInitialLevel,
		Tempers:   NewTempers(),
		PlantedAt: time.Now(),
		SeedMeta:  seed.SeedMeta,
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: plantInteractionRadius,
		},
	}, nil
}

func (p *Plant) Alive() bool {
	return !p.Dead
}

func xpRequiredForLevel(level int64) int64 {
	return int64(math.Round(75 * (math.Pow(float64(level), 2.0) - float64(level))))
}

func (p *Plant) AddXp(xp int64) bool {
	p.Xp += xp

	if nextLevelXpReq := xpRequiredForLevel(p.Level + 1); p.Xp >= nextLevelXpReq {
		p.Xp -= nextLevelXpReq
		p.LevelUp()
		return true
	}

	return false
}

func (p *Plant) LevelUp() int64 {
	p.Level += 1
	return p.Level
}

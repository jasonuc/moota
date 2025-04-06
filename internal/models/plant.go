package models

import (
	"fmt"
	"time"
)

const PLANT_INTERACTION_RADIUSM = 2.0 // TODO: This value is up for debate

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
	circleMeta := CircleMeta{radiusM: PLANT_INTERACTION_RADIUSM, centre: centre}
	if !IsInsideSoil(circleMeta, soil) {
		return nil, fmt.Errorf("plant is not completely inside soil")
	}

	healthBonus := 0.0
	if seed.SeedMeta.OptimalSoil == soil.Type {
		healthBonus += 1.5
	} else {
		if seed.SeedMeta.IsCompatibleWithSoil(soil.Type) {
			healthBonus += 0.5
		} else {
			healthBonus -= 0.5
		}
	}

	// TODO: Compe up with an Xp and Levelling system
	// TODO: Work on giving plants bonus xp when planted on optimal soil

	return &Plant{
		Nickname:  nickname,
		Health:    seed.Health + healthBonus,
		Soil:      soil,
		Tempers:   NewTempers(),
		PlantedAt: time.Now(),
		SeedMeta:  seed.SeedMeta,
		CircleMeta: CircleMeta{
			centre:  centre,
			radiusM: PLANT_INTERACTION_RADIUSM,
		},
	}, nil
}

func (p *Plant) Alive() bool {
	return !p.Dead
}

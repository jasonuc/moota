package models

import (
	"errors"
	"math"
	"time"
)

const (
	PlantInteractionRadius = 3.0
	wateringPlantXpGain    = 30 // TODO: Up for debate
	wateringPlantHpGain    = 5
	minWateringInterval    = 6 * time.Hour
)

var (
	ErrPlantInCooldown     = errors.New("plant in cooldown mode")
	ErrPlantNotFullyInSoil = errors.New("plant not fully inside soil")
	ErrPlantNotFound       = errors.New("plant not found")
)

type PlantAction int

const (
	PlantActionWater PlantAction = iota
)

func ValidPlantAction(action int) bool {
	m := map[PlantAction]bool{
		PlantActionWater: true,
	}
	return m[PlantAction(action)]
}

type Plant struct {
	ID              string    `json:"id"`
	Nickname        string    `json:"nickname"`
	Hp              float64   `json:"hp"`
	Dead            bool      `json:"dead"`
	Activated       bool      `json:"activated"`
	OwnerID         string    `json:"ownerID"`
	Soil            *Soil     `json:"soil,omitempty"`
	Tempers         *Tempers  `json:"tempers,omitempty"`
	TimePlanted     time.Time `json:"timePlanted"`
	LastWateredTime time.Time `json:"lastWateredAt"`
	LastActionTime  time.Time `json:"lastActionTime"`
	SeedMeta
	LevelMeta
	CircleMeta
}

type PlantWithDistanceMFromUser struct {
	Plant
	DistanceM float64 `json:"distanceM"`
}

func NewPlant(seed *Seed, soil *Soil, centre Coordinates) (*Plant, error) {
	if seed.Planted {
		return nil, ErrSeedAlreadyPlanted
	}

	nickname := "Atura" // TODO: Make a function to generate random whimsical names
	seed.Planted = true
	circleMeta := CircleMeta{radiusM: PlantInteractionRadius, C: centre}
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
		Hp:        seed.Hp + healthOffset,
		Soil:      soil,
		OwnerID:   seed.OwnerID,
		Dead:      false,
		Activated: false,
		LevelMeta: NewLeveLMeta(1, xpBonus),
		Tempers:   NewTempers(),
		SeedMeta:  seed.SeedMeta,
		CircleMeta: CircleMeta{
			C:       centre,
			radiusM: PlantInteractionRadius,
		},
	}, nil
}

func (p *Plant) Action(action PlantAction, t time.Time) (bool, error) {
	// TODO: consider converting all time to UTC before carrying out any action
	p.preActionHook(t)

	if !p.Alive() {
		return p.Alive(), nil
	}

	switch action {
	case PlantActionWater:
		// TODO: move this to it's own function and make use of the Soil.WaterRetention
		if t.Sub(p.LastActionTime) > minWateringInterval {
			p.addXp(wateringPlantXpGain)
			p.changeHp(wateringPlantHpGain)
			p.LastWateredTime = t
		} else {
			return p.Alive(), ErrPlantInCooldown
		}
	}
	p.LastActionTime = t

	return p.Alive(), nil
}

func (p *Plant) Refresh(t time.Time) bool {
	p.preActionHook(t)
	return p.Alive()
}

func (p *Plant) preActionHook(t time.Time) {
	// Hp reduction for plant neglect
	hoursSinceLastAction := t.Sub(p.LastActionTime).Hours()
	decreaseMult := math.Floor(hoursSinceLastAction - 12)
	if decreaseMult > 0 {
		alive := p.changeHp(-5 * decreaseMult)
		if !alive {
			return
		}
	}

	// Hp reduction for lack of watering
	hoursSinceLastWatering := math.Floor(t.Sub(p.LastWateredTime).Hours())
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

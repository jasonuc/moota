package models

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

const (
	PlantInteractionRadius = 3.0
	wateringPlantXpGain    = 30
	wateringPlantHpGain    = 5
	minWateringInterval    = 6 * time.Hour

	hpPenaltyPlantNeglect = 0.1
	hpPenaltyLackOfWater  = 0.3
)

var (
	ErrPlantInCooldown     = errors.New("plant is in cooldown mode")
	ErrPlantNotFullyInSoil = errors.New("plant not fully inside soil")
	ErrPlantNotFound       = errors.New("plant not found")
)

type PlantAction int

const (
	PlantActionWater PlantAction = iota + 1
)

func ValidPlantAction(action int) bool {
	m := map[PlantAction]bool{
		PlantActionWater: true,
	}
	return m[PlantAction(action)]
}

type Plant struct {
	ID              string     `json:"id"`
	Nickname        string     `json:"nickname"`
	Hp              float64    `json:"hp"`
	Dead            bool       `json:"dead"`
	OwnerID         string     `json:"ownerID"`
	Soil            *Soil      `json:"soil,omitempty"`
	Tempers         *Tempers   `json:"tempers,omitempty"`
	TimePlanted     time.Time  `json:"timePlanted"`
	TimeOfDeath     *time.Time `json:"timeOfDeath"`
	LastWateredTime time.Time  `json:"lastWateredAt"`
	LastActionTime  time.Time  `json:"lastActionTime"`
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

	nickname := generateNickname()
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
	p.preActionHook(t)

	if !p.Alive() {
		return p.Alive(), nil
	}

	switch action {
	case PlantActionWater:
		if t.Sub(p.LastActionTime) > minWateringInterval || p.TimePlanted.Equal(p.LastActionTime) {
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
	decreaseMult := math.Floor(hoursSinceLastAction - 24)
	if decreaseMult > 0 {
		alive := p.changeHp(-hpPenaltyPlantNeglect * decreaseMult)
		if !alive {
			return
		}
	}

	// Hp reduction for lack of watering
	hoursSinceLastWatered := math.Floor(t.Sub(p.LastWateredTime).Hours())
	decreaseMult = math.Floor(hoursSinceLastWatered - 12)
	if decreaseMult > 0 {
		p.changeHp(-hpPenaltyLackOfWater * decreaseMult)
	}
}

func (p *Plant) Alive() bool {
	return !p.Dead
}

func (p *Plant) changeHp(delta float64) bool {
	p.Hp = math.Max(0, math.Min(100, p.Hp+delta)) // clamp hp between 0 and 100
	if p.Hp == 0 {
		p.Die(time.Now())
	}
	return p.Alive()
}

func (p *Plant) Die(timeOfDeath time.Time) {
	p.Hp = 0
	p.Dead = true
	p.TimeOfDeath = &timeOfDeath
}

func generateNickname() string {
	adjectives := []string{
		"Wiggly", "Sparkle", "Fuzzy", "Giggly", "Dapper", "Sneaky", "Wobble",
		"Whispering", "Disco", "Dazzle", "Twinkle", "Bubbly", "Misty", "Cosmic",
		"Sleepy", "Grumpy", "Zany", "Mischievous", "Bouncy", "Breezy", "Snazzy",
		"Quirky", "Mystical", "Rainbow", "Doodle", "Fluffy", "Zigzag", "Velvet",
	}

	nouns := []string{
		"Sprout", "Leaf", "Blossom", "Twig", "Root", "Petal", "Bean", "Shadow",
		"Whisker", "Pickle", "Pancake", "Buttons", "Socks", "Teapot", "Muffin",
		"Rocket", "Banjo", "Noodle", "Pebble", "Bubble", "Jelly", "Whistle",
		"Wizard", "Moonbeam", "Droplet", "Thunder", "Feather", "Crumble",
	}

	titles := []string{
		"the Magnificent", "the Curious", "the Bold", "the Shy", "the Mighty",
		"the Daring", "the Gentle", "the Wise", "the Fabulous", "the Ticklish",
		"the Brave", "the Mysterious", "the Fancy", "the Dreamer", "the Champion",
		"the Explorer", "the Melodious", "the Inventor", "von Bloomenstein",
		"the Magical", "the Whimsical", "the Chaotic", "Jr.", "the Majestic",
	}

	adjective := adjectives[rand.IntN(len(adjectives))]
	noun := nouns[rand.IntN(len(nouns))]

	if rand.IntN(2) == 0 {
		title := titles[rand.IntN(len(titles))]
		return fmt.Sprintf("%s %s %s", adjective, noun, title)
	}

	return fmt.Sprintf("%s %s", adjective, noun)
}

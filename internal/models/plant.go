package models

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

const (
	PlantInteractionRadius = 20 // TODO: This value is still experimental
	wateringPlantXpGain    = 30
	wateringPlantHpGain    = 5

	hpDecayInterval = 4 * time.Hour

	wateringGracePeriod = 4 * time.Hour
	wateringCooldown    = 3 * time.Hour

	minRefreshInterval = 5 * time.Minute
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
	ID                string     `json:"id"`
	Nickname          string     `json:"nickname"`
	Hp                float64    `json:"hp"`
	Dead              bool       `json:"dead"`
	OwnerID           string     `json:"ownerID"`
	Soil              *Soil      `json:"soil,omitempty"`
	Tempers           *Tempers   `json:"tempers,omitempty"`
	TimePlanted       time.Time  `json:"timePlanted"`
	TimeOfDeath       *time.Time `json:"timeOfDeath"`
	LastWateredAt     time.Time  `json:"lastWateredAt"`
	LastActionAt      time.Time  `json:"lastActionAt"`
	LastRefreshedAt   *time.Time `json:"lastRefreshedAt,omitempty"`
	GracePeriodEndsAt *time.Time `json:"gracePeriodEndsAt,omitempty"`
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
		LastRefreshedAt:   nil,
		GracePeriodEndsAt: nil,
	}, nil
}

func (p *Plant) Action(action PlantAction, t time.Time) (bool, error) {
	p.applyTimeBasedChanges(t)

	if !p.Alive() {
		return p.Alive(), nil
	}

	switch action {
	case PlantActionWater:
		if p.LastWateredAt.IsZero() || t.Sub(p.LastWateredAt) >= wateringCooldown {
			p.addXp(wateringPlantXpGain)
			p.changeHp(wateringPlantHpGain)
			p.LastWateredAt = t

			gracePeriodEnd := t.Add(wateringGracePeriod)
			p.GracePeriodEndsAt = &gracePeriodEnd
		} else {
			return p.Alive(), ErrPlantInCooldown
		}
	}

	p.LastActionAt = t
	return p.Alive(), nil
}

func (p *Plant) Refresh(t time.Time) bool {
	if p.LastRefreshedAt != nil && t.Sub(*p.LastRefreshedAt) < minRefreshInterval {
		return p.Alive()
	}

	p.applyTimeBasedChanges(t)
	return p.Alive()
}

func (p *Plant) applyTimeBasedChanges(t time.Time) {
	if p.Dead {
		return
	}

	var lastCalculatedAt time.Time
	if p.LastRefreshedAt != nil {
		lastCalculatedAt = *p.LastRefreshedAt
	} else {
		lastCalculatedAt = p.TimePlanted
	}

	p.LastRefreshedAt = &t

	p.calculateAndApplyDecay(lastCalculatedAt, t)
}

func (p *Plant) calculateAndApplyDecay(fromTime, toTime time.Time) {
	if p.Dead {
		return
	}

	elapsed := toTime.Sub(fromTime)
	totalIntervals := int(elapsed / hpDecayInterval)

	if totalIntervals == 0 {
		return
	}

	intervalsToApply := 0

	for i := 0; i < totalIntervals; i++ {
		intervalTime := fromTime.Add(time.Duration(i+1) * hpDecayInterval)

		inGracePeriod := false
		if p.GracePeriodEndsAt != nil && intervalTime.Before(*p.GracePeriodEndsAt) {
			inGracePeriod = true
		}

		if !inGracePeriod {
			intervalsToApply++
		}
	}

	if intervalsToApply > 0 {
		p.changeHp(-float64(intervalsToApply))
	}
}

func (p *Plant) IsInGracePeriod(t time.Time) bool {
	return p.GracePeriodEndsAt != nil && t.Before(*p.GracePeriodEndsAt)
}

func (p *Plant) CanBeWatered(t time.Time) bool {
	if p.Dead {
		return false
	}
	return p.LastWateredAt.IsZero() || t.Sub(p.LastWateredAt) >= wateringCooldown
}

func (p *Plant) TimeUntilNextWatering(t time.Time) time.Duration {
	if p.CanBeWatered(t) {
		return 0
	}
	return wateringCooldown - t.Sub(p.LastWateredAt)
}

func (p *Plant) TimeUntilGracePeriodEnds(t time.Time) time.Duration {
	if !p.IsInGracePeriod(t) {
		return 0
	}
	return p.GracePeriodEndsAt.Sub(t)
}

func (p *Plant) Alive() bool {
	return !p.Dead
}

func (p *Plant) changeHp(delta float64) bool {
	p.Hp = math.Max(0, math.Min(100, p.Hp+delta))
	if p.Hp == 0 {
		p.Die(time.Now())
	}
	return p.Alive()
}

func (p *Plant) Die(timeOfDeath time.Time) {
	p.Hp = 0
	p.Dead = true
	p.TimeOfDeath = &timeOfDeath
	p.GracePeriodEndsAt = nil
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

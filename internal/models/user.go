package models

import (
	"errors"
	"sort"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Title        string    `json:"title"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	LevelMeta
}

type UserProfile struct {
	Username        string   `json:"username"`
	Title           string   `json:"title"`
	Level           int64    `json:"level"`
	Top3AlivePlants []*Plant `json:"top3AlivePlants"`
	DeceasedPlants  []*Plant `json:"deceasedPlants"`
	PlantCount      `json:"plantCount"`
	SeedCount       `json:"seedCount"`
}

type PlantCount struct {
	Alive    int64 `json:"alive"`
	Deceased int64 `json:"deceased"`
}

type SeedCount struct {
	Planted int64 `json:"planted"`
	Unused  int64 `json:"unused"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewUserProfile(user *User, plantCount *PlantCount, seedCount *SeedCount, plants []*Plant) *UserProfile {
	userProfile := new(UserProfile)

	userProfile.Username = user.Username
	userProfile.Title = user.Title
	userProfile.Level = user.Level

	userProfile.PlantCount = *plantCount
	userProfile.SeedCount = *seedCount

	userProfile.Top3AlivePlants = make([]*Plant, 0)
	userProfile.DeceasedPlants = make([]*Plant, 0)

	alivePlants := make([]*Plant, 0)
	deceasedPlants := make([]*Plant, 0)

	for _, plant := range plants {
		plant.CircleMeta = CircleMeta{}

		if plant.Dead {
			deceasedPlants = append(deceasedPlants, plant)
		} else {
			alivePlants = append(alivePlants, plant)
		}
	}

	sort.Slice(alivePlants, func(i, j int) bool {
		return alivePlants[i].Hp > alivePlants[j].Hp
	})

	maxAlive := min(len(alivePlants), 3)
	userProfile.Top3AlivePlants = alivePlants[:maxAlive]

	sort.Slice(deceasedPlants, func(i, j int) bool {
		if deceasedPlants[i].TimeOfDeath == nil || deceasedPlants[j].TimeOfDeath == nil {
			return deceasedPlants[i].TimeOfDeath != nil
		}
		return deceasedPlants[i].TimeOfDeath.After(*deceasedPlants[j].TimeOfDeath)
	})

	userProfile.DeceasedPlants = deceasedPlants

	return userProfile
}

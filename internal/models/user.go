package models

import (
	"errors"
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
	Username   string `json:"username"`
	Title      string `json:"title"`
	Level      int64  `json:"level"`
	PlantCount `json:"plantCount"`
	SeedCount  `json:"seedCount"`
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

func NewUserProfile(user *User, plantCount *PlantCount, seedCount *SeedCount) *UserProfile {
	userProfile := new(UserProfile)

	userProfile.Username = user.Username
	userProfile.Title = user.Title
	userProfile.Level = user.Level

	userProfile.PlantCount = *plantCount
	userProfile.SeedCount = *seedCount

	return userProfile
}

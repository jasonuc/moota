package models

import "time"

type User struct {
	ID           string
	Username     string
	Title        string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LevelMeta
}

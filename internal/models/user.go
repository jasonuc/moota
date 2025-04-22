package models

import "time"

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

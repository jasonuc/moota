package models

import (
	"errors"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshToken struct {
	ID        string
	UserID    string
	Hash      []byte
	Plain     string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

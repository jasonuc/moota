package models

import (
	"errors"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
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

package models

import "time"

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`

	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
	Issuer    string    `json:"iss"`
}

type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userID"`
	Hash      []byte     `json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	ExpiresAt time.Time  `json:"expiresAt"`
	RevokedAt *time.Time `json:"revokedAt,omitempty"`
}

package store

import (
	"github.com/jasonuc/moota/internal/models"
)

type RefreshTokenStore interface {
	Insert(*models.RefreshToken) error
	GetByHash([]byte) (*models.RefreshToken, error)
	Revoke(string) error
}

type refreshTokenStore struct {
	db Querier
}

func (s *refreshTokenStore) Insert(refreshToken *models.RefreshToken) error {
	return nil
}

func (s *refreshTokenStore) GetByHash(refreshTokenHash []byte) (*models.RefreshToken, error) {
	return nil, nil
}

func (s *refreshTokenStore) Revoke(id string) error {
	return nil
}

package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type seedStore struct {
	db *sql.DB
}

func (s *seedStore) Get(id string) (*models.Seed, error) {
	return nil, nil
}

func (s *seedStore) Insert(*models.Seed) error {
	return nil
}

func (s *seedStore) Delete(id string) error {
	return nil
}

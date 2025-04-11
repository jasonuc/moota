package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type soilStore struct {
	db *sql.DB
}

func (s *soilStore) Get(id string) (*models.Soil, error) {
	return nil, nil
}

func (s *soilStore) Insert(soil *models.Soil) error {
	return nil
}

func (s *soilStore) Delete(id string) error {
	return nil
}

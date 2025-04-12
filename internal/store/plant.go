package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type plantStore struct {
	db *sql.DB
}

func (s *plantStore) Get(id string) (*models.Plant, error) {
	return nil, nil
}

func (s *plantStore) Insert(plant *models.Plant) error {
	return nil
}

func (s *plantStore) Update(updatedPlant *models.Plant) error {
	return nil
}

func (s *plantStore) Delete(id string) error {
	return nil
}

package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type plantStore struct {
	db *sql.DB
}

func (p *plantStore) Get(id string) (*models.Plant, error) {
	return nil, nil
}

func (p *plantStore) Insert(plant *models.Plant) error {
	return nil
}

func (p *plantStore) Update(updatedPlant *models.Plant) error {
	return nil
}

func (p *plantStore) Delete(id string) error {
	return nil
}

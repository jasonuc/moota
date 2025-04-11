package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type Store struct {
	Plant interface {
		Get(string) (*models.Plant, error)
		Insert(*models.Plant) error
		Update(*models.Plant) error
		Delete(string) error
	}

	Soil interface {
		Get(string) (*models.Soil, error)
		Insert(*models.Soil) error
		Delete(string) error
	}

	Seed interface {
		Get(string) (*models.Seed, error)
		Insert(*models.Seed) error
		Delete(string) error
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Seed:  &seedStore{db},
		Plant: &plantStore{db},
		Soil:  &soilStore{db},
	}
}

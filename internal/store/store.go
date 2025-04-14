package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type Store struct {
	User interface {
		Insert(*models.User) error
		GetByEmail(string) (*models.User, error)
		GetByID(string) (*models.User, error)
		GetByUsername(string) (*models.User, error)
		Update(*models.User) error
		Delete(string) error
	}

	Plant interface {
		Get(string) (*models.Plant, error)
		GetAllByOwnerID(string) ([]*models.Plant, error)
		GetAllNearCoordinates(models.Coordinates, float64) ([]*models.Plant, error)
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
		GetAllByOwnerID(string) ([]*models.Seed, error)
		Insert(*models.Seed) error
		MarkAsPlanted(*models.Seed) error
		Delete(string) error
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		User:  &userStore{db},
		Seed:  &seedStore{db},
		Plant: &plantStore{db},
		Soil:  &soilStore{db},
	}
}

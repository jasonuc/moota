package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jasonuc/moota/internal/models"
)

type Store struct {
	db *sql.DB

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
		GetAllInSoilAndInProximity(string, models.Coordinates, float64) ([]*models.Plant, error)
		GetByOwnerIDAndProximity(string, models.Coordinates, float64) ([]*models.Plant, error)
		GetByOwnerIDAndOrderByProximity(string, models.Coordinates) ([]*models.Plant, error)
		ActivatePlant(string) error
		Insert(*models.Plant) error
		Update(*models.Plant) error
		Delete(string) error
	}

	Soil interface {
		Get(string) (*models.Soil, error)
		GetAllInProximity(models.Coordinates, float64) ([]*models.Soil, error)
		Insert(*models.Soil) error
		Delete(string) error
	}

	Seed interface {
		Get(string) (*models.Seed, error)
		GetAllByOwnerID(string) ([]*models.Seed, error)
		Insert(*models.Seed) error
		MarkAsPlanted(string) error
		Delete(string) error
	}
}

var (
	ErrTransactionCouldNotStart = errors.New("transaction could not be started")
)

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:    db,
		User:  &userStore{db},
		Seed:  &seedStore{db},
		Plant: &plantStore{db},
		Soil:  &soilStore{db},
	}
}

func (s *Store) Begin() (*Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

func (s *Store) WithTx(tx *Transaction) *Store {
	return &Store{
		User:  &userStore{tx.tx},
		Seed:  &seedStore{tx.tx},
		Plant: &plantStore{tx.tx},
		Soil:  &soilStore{tx.tx},
	}
}

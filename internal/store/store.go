package store

import (
	"context"
	"database/sql"
	"errors"
)

type Store struct {
	db    *sql.DB
	User  UserStore
	Plant PlantStore
	Soil  SoilStore
	Seed  SeedStore
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

func (s *Store) WithTx(transaction *Transaction) *Store {
	return &Store{
		User:  &userStore{transaction.tx},
		Seed:  &seedStore{transaction.tx},
		Plant: &plantStore{transaction.tx},
		Soil:  &soilStore{transaction.tx},
	}
}

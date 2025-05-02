package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jasonuc/moota/internal/models"
)

type SeedStore interface {
	Get(context.Context, string) (*models.Seed, error)
	GetAllByOwnerID(context.Context, string) ([]*models.Seed, error)
	GetCountByUsername(context.Context, string) (*models.SeedCount, error)
	Insert(context.Context, *models.Seed) error
	MarkAsPlanted(context.Context, string) error
	Delete(context.Context, string) error
}

type seedStore struct {
	db Querier
}

func (s *seedStore) GetCountByUsername(ctx context.Context, username string) (*models.SeedCount, error) {
	q := `SELECT s.planted, count(*) FROM seeds s
			JOIN users u ON s.owner_id = u.id
			WHERE u.username = $1
			GROUP BY s.planted ORDER BY s.planted ASC;`

	rows, err := s.db.QueryContext(ctx, q, username)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	seedCount := &models.SeedCount{}

	for rows.Next() {
		var planted bool
		var count int64

		err := rows.Scan(&planted, &count)
		if err != nil {
			return nil, err
		}

		if planted {
			seedCount.Planted = count
		} else {
			seedCount.Unused = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seedCount, nil
}

func (s *seedStore) GetAllByOwnerID(ctx context.Context, ownerID string) ([]*models.Seed, error) {
	q := `SELECT id, owner_id, hp, planted, optimal_soil, botanical_name, created_at FROM seeds
			WHERE owner_id = $1 AND planted = false;`

	rows, err := s.db.QueryContext(ctx, q, ownerID)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	seeds := make([]*models.Seed, 0)

	for rows.Next() {
		seed := new(models.Seed)
		err := rows.Scan(&seed.ID, &seed.OwnerID, &seed.Hp, &seed.Planted, &seed.OptimalSoil, &seed.BotanicalName, &seed.CreatedAt)
		if err != nil {
			return nil, err
		}

		seeds = append(seeds, seed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seeds, nil
}

func (s *seedStore) Get(ctx context.Context, id string) (*models.Seed, error) {
	q := `SELECT id, owner_id, hp, planted, optimal_soil, botanical_name, created_at FROM seeds
			WHERE id = $1;`

	seed := new(models.Seed)
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&seed.ID, &seed.OwnerID, &seed.Hp, &seed.Planted, &seed.OptimalSoil, &seed.BotanicalName, &seed.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, models.ErrSeedNotFound
		default:
			return nil, err
		}
	}

	return seed, nil
}

func (s *seedStore) Insert(ctx context.Context, seed *models.Seed) error {
	q := `INSERT INTO seeds (owner_id, hp, planted, optimal_soil, botanical_name)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at;`

	err := s.db.QueryRowContext(ctx, q, seed.OwnerID, seed.Hp, seed.Planted, seed.OptimalSoil, seed.BotanicalName).Scan(
		&seed.ID, &seed.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *seedStore) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM seeds 
			WHERE id = $1;`

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrSeedNotFound
	}

	return nil
}

func (s *seedStore) MarkAsPlanted(ctx context.Context, seedID string) error {
	q := `UPDATE seeds
		SET planted = True
		WHERE id = $1;`

	res, err := s.db.ExecContext(ctx, q, seedID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrSeedNotFound
	}

	return nil
}

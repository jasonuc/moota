package store

import (
	"github.com/jasonuc/moota/internal/models"
)

type seedStore struct {
	db dbOrTx
}

func (s *seedStore) GetAllByOwnerID(ownerID string) ([]*models.Seed, error) {
	q := `SELECT id, owner_id, health, planted, optimal_soil, botanical_name, created_at FROM seeds
			WHERE owner_id = $1 AND planted = false;`

	rows, err := s.db.Query(q, ownerID)
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

func (s *seedStore) Get(id string) (*models.Seed, error) {
	q := `SELECT id, owner_id, health, planted, optimal_soil, botanical_name, created_at FROM seeds
			WHERE id = $1;`

	seed := new(models.Seed)
	err := s.db.QueryRow(q, id).Scan(
		&seed.ID, &seed.OwnerID, &seed.Hp, &seed.Planted, &seed.OptimalSoil, &seed.BotanicalName, &seed.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return seed, nil
}

func (s *seedStore) Insert(seed *models.Seed) error {
	q := `INSERT INTO seeds (owner_id, health, planted, optimal_soil, botanical_name)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at;`

	err := s.db.QueryRow(q, seed.OwnerID, seed.Hp, seed.Planted, seed.OptimalSoil, seed.BotanicalName).Scan(
		&seed.ID, &seed.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *seedStore) Delete(id string) error {
	q := `DELETE FROM seeds 
			WHERE id = $1;`

	res, err := s.db.Exec(q, id)
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

func (s *seedStore) MarkAsPlanted(seed *models.Seed) error {
	if !seed.Planted {
		return models.ErrSeedAlreadyPlanted
	}

	q := `UPDATE seeds
		SET planted = $2
		WHERE id = $1;`

	res, err := s.db.Exec(q, seed.ID, seed.Planted)
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

package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jasonuc/moota/internal/models"
)

type UserStore interface {
	Insert(context.Context, *models.User) error
	GetByEmail(context.Context, string) (*models.User, error)
	GetByID(context.Context, string) (*models.User, error)
	GetByUsername(context.Context, string) (*models.User, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, string) error
}

type userStore struct {
	db Querier
}

func (s *userStore) Insert(ctx context.Context, user *models.User) error {
	q := `INSERT INTO users (username, email, password_hash, created_at, updated_at, level, xp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at;`

	err := s.db.QueryRowContext(
		ctx, q, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.Level, user.Xp,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *userStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at, level, xp
		FROM users WHERE email = $1;`

	user := &models.User{}
	err := s.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Level, &user.Xp,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userStore) GetByID(ctx context.Context, id string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at, level, xp
		FROM users WHERE id = $1;`

	user := &models.User{}
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Level, &user.Xp,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userStore) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at, level, xp
		FROM users WHERE username = $1;`

	user := &models.User{}
	err := s.db.QueryRowContext(ctx, q, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Level, &user.Xp,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userStore) Update(ctx context.Context, updatedUser *models.User) error {
	q := `UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = $4, level = $5, xp = $6
		WHERE id = $7;`

	res, err := s.db.ExecContext(ctx, q,
		updatedUser.Username, updatedUser.Email, updatedUser.PasswordHash,
		updatedUser.UpdatedAt, updatedUser.Level, updatedUser.Xp, updatedUser.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

func (s *userStore) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM users WHERE id = $1;`

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

package store

import (
	"database/sql"
	"errors"

	"github.com/jasonuc/moota/internal/models"
)

type UserStore interface {
	Insert(*models.User) error
	GetByEmail(string) (*models.User, error)
	GetByID(string) (*models.User, error)
	GetByUsername(string) (*models.User, error)
	Update(*models.User) error
	Delete(string) error
}

type userStore struct {
	db Querier
}

func (s *userStore) Insert(user *models.User) error {
	q := `INSERT INTO users (id, username, email, password_hash, created_at, updated_at, level, xp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at;`

	err := s.db.QueryRow(
		q, user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.Level, user.Xp,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *userStore) GetByEmail(email string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at, level, xp
		FROM users WHERE email = $1;`

	user := &models.User{}
	err := s.db.QueryRow(q, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Level, &user.Xp,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *userStore) GetByID(id string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at, level, xp
		FROM users WHERE id = $1;`

	user := &models.User{}
	err := s.db.QueryRow(q, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Level, &user.Xp,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *userStore) GetByUsername(username string) (*models.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at, level, xp
		FROM users WHERE username = $1;`

	user := &models.User{}
	err := s.db.QueryRow(q, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.Level, &user.Xp,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *userStore) Update(updatedUser *models.User) error {
	q := `UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = $4, level = $5, xp = $6
		WHERE id = $7;`

	res, err := s.db.Exec(q,
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

func (s *userStore) Delete(id string) error {
	q := `DELETE FROM users WHERE id = $1;`

	res, err := s.db.Exec(q, id)
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

package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type userStore struct {
	db *sql.DB
}

func (s *userStore) Insert(user *models.User) error {
	return nil
}

func (s *userStore) GetByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (s *userStore) GetByID(id string) (*models.User, error) {
	return nil, nil
}

func (s *userStore) GetByUsername(username string) (*models.User, error) {
	return nil, nil
}

func (s *userStore) Update(updatedUser *models.User) error {
	return nil
}

func (s *userStore) Delete(id string) error {
	return nil
}

package store

import (
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

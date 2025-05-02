package services

import (
	"context"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type UserService interface {
	GetUserByID(context.Context, string) (*models.User, error)
}

type userService struct {
	store *store.Store
}

func NewUserService(store *store.Store) UserService {
	return &userService{
		store: store,
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.store.User.GetByID(ctx, userID)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	return user, nil
}

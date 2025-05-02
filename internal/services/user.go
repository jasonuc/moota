package services

import (
	"context"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type UserService interface {
	GetUser(context.Context, string) (*models.User, error)
	GetUserProfile(context.Context, string) (*models.UserProfile, error)
}

type userService struct {
	store *store.Store
}

func NewUserService(store *store.Store) UserService {
	return &userService{
		store: store,
	}
}

func (s *userService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.store.User.GetByID(ctx, userID)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	return user, nil
}

func (s *userService) GetUserProfile(ctx context.Context, username string) (*models.UserProfile, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	user, err := tx.User.GetByUsername(ctx, username)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	plantCount, err := tx.Plant.GetCountByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	seedCount, err := tx.Seed.GetCountByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	userProfile := models.NewUserProfile(user, plantCount, seedCount)

	return userProfile, nil
}

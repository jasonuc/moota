package services

import (
	"context"
	"fmt"
	"unicode"

	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/dto"
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

func (s *userService) UpdateUser(ctx context.Context, userID string, dto dto.UpdateUserReq) (*models.User, error) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if userID != userIDFromCtx {
		return nil, fmt.Errorf("you do not have authorised access")
	}

	user, err := s.GetUserByID(ctx, userIDFromCtx)
	if err != nil {
		return nil, err
	}

	if !isValidUsername(dto.Username) {
		return nil, ErrInvalidUsername
	}

	user.Username = dto.Username
	if err := s.store.User.Update(ctx, user); err != nil {
		switch {
		default:
			return nil, err
		}
	}

	return user, nil
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 8 {
		return false
	}

	for _, char := range username {
		if !unicode.IsLetter(char) {
			return false
		}
	}

	return true
}

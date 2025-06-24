package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidUsername       = errors.New("invalid username")
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrTokenExpiredOrRevoked = errors.New("token expired or revoked")
)

type AuthService interface {
	Register(context.Context, dto.UserRegisterReq) (*models.User, *models.TokenPair, error)
	Login(context.Context, dto.UserLoginReq) (*models.TokenPair, error)
	RefreshTokens(context.Context, string) (*models.TokenPair, error)
	VerifyAccessToken(context.Context, string) (string, error)
	ChangeUserUsername(context.Context, string, dto.ChangeUsernameReq) (*models.User, error)
	ChangeUserEmail(context.Context, string, dto.ChangeEmailReq) (*models.User, error)
	ChangeUserPassword(context.Context, string, dto.ChangePasswordReq) (*models.User, error)
	GetAccessTokenTTL() int
	GetRefreshTokenTTL() int
}

type authService struct {
	store           *store.Store
	accessSecret    []byte
	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
	issuer          string
}

func NewAuthService(store *store.Store, acessSecret []byte, refreshTTL, acessTTL time.Duration, issuer string) AuthService {
	return &authService{
		store:           store,
		accessSecret:    acessSecret,
		refreshTokenTTL: refreshTTL,
		accessTokenTTL:  acessTTL,
		issuer:          issuer,
	}
}

func (s *authService) Register(ctx context.Context, dto dto.UserRegisterReq) (*models.User, *models.TokenPair, error) {
	_, err := s.store.User.GetByEmail(ctx, dto.Email)
	if err == nil {
		return nil, nil, ErrInvalidEmail
	}

	_, err = s.store.User.GetByUsername(ctx, dto.Email)
	if err == nil {
		return nil, nil, ErrInvalidUsername
	}

	if !isValidUsername(dto.Username) {
		return nil, nil, ErrInvalidUsername
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: hashedPassword,
		LevelMeta:    models.NewLeveLMeta(1, 0),
	}

	err = s.store.User.Insert(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, nil, err
	}

	if err := s.store.RefreshToken.Insert(ctx, refreshToken); err != nil {
		return nil, nil, err
	}

	tokenPair := &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Plain,
	}

	return user, tokenPair, nil
}

func (s *authService) Login(ctx context.Context, dto dto.UserLoginReq) (*models.TokenPair, error) {
	user, err := s.store.User.GetByEmail(ctx, dto.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	if err := s.store.RefreshToken.Insert(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Plain,
	}, nil
}

func (s *authService) RefreshTokens(ctx context.Context, tokenRefresh string) (*models.TokenPair, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()
	tx := s.store.WithTx(transaction)

	tokenHash := sha256.Sum256([]byte(tokenRefresh))

	token, err := tx.RefreshToken.GetByHash(ctx, tokenHash[:])
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if token.ExpiresAt.Before(time.Now()) || token.RevokedAt != nil {
		return nil, ErrTokenExpiredOrRevoked
	}

	user, err := tx.User.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	if err := tx.RefreshToken.Revoke(ctx, token.ID); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	if err := tx.RefreshToken.Insert(ctx, refreshToken); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Plain,
	}, nil
}

func (s *authService) VerifyAccessToken(ctx context.Context, accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err := claims.GetSubject()
		if err != nil {
			return "", fmt.Errorf("invalid user id in token: %w", err)
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

func (s *authService) ChangeUserUsername(ctx context.Context, userID string, dto dto.ChangeUsernameReq) (*models.User, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	if _, err := tx.User.GetByUsername(ctx, dto.NewUsername); err == nil {
		return nil, ErrInvalidUsername
	}

	user, err := tx.User.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !isValidUsername(dto.NewUsername) {
		return nil, ErrInvalidUsername
	}

	user.Username = dto.NewUsername
	if err := tx.User.Update(ctx, user); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) ChangeUserEmail(ctx context.Context, userID string, dto dto.ChangeEmailReq) (*models.User, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	if _, err := tx.User.GetByEmail(ctx, dto.NewEmail); err == nil {
		return nil, ErrInvalidEmail
	}

	user, err := tx.User.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Email = dto.NewEmail
	if err := tx.User.Update(ctx, user); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) ChangeUserPassword(ctx context.Context, userID string, dto dto.ChangePasswordReq) (*models.User, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	user, err := tx.User.GetByID(ctx, userID)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(dto.OldPassword)); err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, ErrInvalidCredentials
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = newPasswordHash
	if err := tx.User.Update(ctx, user); err != nil {
		return nil, err
	}

	// this current approach to change password would revoke all the user's refresh tokens which in effect signs them out when the refresh token expires
	// the current approach does have a problem though. It is going to weight until the user needs to refresh their access token before they'd be logged out so...
	// TODO: improve maybe return a new token pair
	if err := tx.RefreshToken.RevokeManyByUserID(ctx, userID); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) generateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	accessExp := time.Now().Add(s.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      accessExp.Unix(),
		"iat":      now.Unix(),
		"iss":      s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.accessSecret)
	if err != nil {
		return "", fmt.Errorf("could not generate access token: %w", err)
	}

	return accessToken, nil
}

func (s *authService) GetAccessTokenTTL() int {
	return int(s.accessTokenTTL)
}

func (s *authService) GetRefreshTokenTTL() int {
	return int(s.refreshTokenTTL)
}

func (s *authService) generateRefreshToken(user *models.User) (*models.RefreshToken, error) {
	now := time.Now()
	refreshExp := time.Now().Add(s.refreshTokenTTL)

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("could not generate refresh token %w", err)
	}

	refreshTokenPlain := base64.RawURLEncoding.EncodeToString(refreshTokenBytes)
	refreshTokenHash := sha256.Sum256([]byte(refreshTokenPlain))

	return &models.RefreshToken{
		UserID:    user.ID,
		Hash:      refreshTokenHash[:],
		Plain:     refreshTokenPlain,
		CreatedAt: now,
		ExpiresAt: refreshExp,
	}, nil
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

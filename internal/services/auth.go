package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Register(string, string, string) (*models.User, error)
	Login(string, string) (*models.TokenPair, error)
	RefreshTokens(string) (*models.TokenPair, error)
	VerifyAccessToken(string) (string, error)
}

type authService struct {
	store           *store.Store
	refreshSecret   []byte
	accessSecret    []byte
	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
	issuer          string
}

func NewAuthService(store *store.Store, refreshSecret, acessSecret []byte, refreshTTL, acessTTL time.Duration, issuer string) AuthService {
	return &authService{
		store:           store,
		refreshSecret:   refreshSecret,
		accessSecret:    acessSecret,
		refreshTokenTTL: refreshTTL,
		accessTokenTTL:  acessTTL,
		issuer:          issuer,
	}
}

func (s *authService) Register(email, username, password string) (*models.User, error) {
	_, err := s.store.User.GetByEmail(email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	_, err = s.store.User.GetByUsername(username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		LevelMeta:    models.NewLeveLMeta(1, 0),
	}

	err = s.store.User.Insert(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(usernameOrEmail, password string) (*models.TokenPair, error) {
	var user *models.User
	var err error

	user, err = s.store.User.GetByEmail(usernameOrEmail)
	if err != nil {
		user, err = s.store.User.GetByUsername(usernameOrEmail)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokenPair, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	if err := s.store.RefreshToken.Insert(refreshToken); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) RefreshTokens(refreshTokenString string) (*models.TokenPair, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()
	tx := s.store.WithTx(transaction)

	tokenHash := sha256.Sum256([]byte(refreshTokenString))

	token, err := tx.RefreshToken.GetByHash(tokenHash[:])
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if token.ExpiresAt.Before(time.Now()) || token.RevokedAt != nil {
		return nil, fmt.Errorf("token expired or revoked")
	}

	user, err := tx.User.GetByID(token.UserID)
	if err != nil {
		return nil, err
	}

	if err := tx.RefreshToken.Revoke(token.ID); err != nil {
		return nil, err
	}

	tokenPair, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	if err := tx.RefreshToken.Insert(refreshToken); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) VerifyAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
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

func (s *authService) generateTokenPair(user *models.User) (*models.TokenPair, *models.RefreshToken, error) {
	now := time.Now()
	accessExp := now.Add(s.accessTokenTTL)
	refreshExp := now.Add(s.refreshTokenTTL)

	acessClaims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      accessExp.Unix(),
		"iss":      s.issuer,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, acessClaims)
	acessTokenString, err := accessToken.SignedString(s.accessSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate acess token: %w", err)
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, nil, fmt.Errorf("could not generate refresh token: %w", err)
	}
	refreshTokenString := base64.RawURLEncoding.EncodeToString(refreshTokenBytes)
	refreshTokenHash := sha256.Sum256([]byte(refreshTokenString))

	return &models.TokenPair{
			AccessToken:  acessTokenString,
			RefreshToken: refreshTokenString,
		}, &models.RefreshToken{
			UserID:    user.ID,
			Hash:      refreshTokenHash[:],
			ExpiresAt: refreshExp,
		}, nil
}

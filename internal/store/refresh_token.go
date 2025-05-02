package store

import (
	"context"

	"github.com/jasonuc/moota/internal/models"
)

type RefreshTokenStore interface {
	Insert(context.Context, *models.RefreshToken) error
	GetByHash(context.Context, []byte) (*models.RefreshToken, error)
	Revoke(context.Context, string) error
	RevokeAllByUserID(context.Context, string) error
}

type refreshTokenStore struct {
	db Querier
}

func (s *refreshTokenStore) Insert(ctx context.Context, refreshToken *models.RefreshToken) error {
	q := `INSERT INTO refresh_tokens (user_id, hash, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, revoked_at;`

	err := s.db.QueryRowContext(
		ctx, q, refreshToken.UserID, refreshToken.Hash, refreshToken.CreatedAt, refreshToken.ExpiresAt,
	).Scan(&refreshToken.ID, &refreshToken.RevokedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *refreshTokenStore) GetByHash(ctx context.Context, refreshTokenHash []byte) (*models.RefreshToken, error) {
	q := `SELECT id, user_id, hash, created_at, expires_at, revoked_at
		FROM refresh_tokens WHERE hash = $1;`

	refreshToken := new(models.RefreshToken)
	err := s.db.QueryRowContext(ctx, q, refreshTokenHash).Scan(
		&refreshToken.ID, &refreshToken.UserID, &refreshToken.Hash,
		&refreshToken.CreatedAt, &refreshToken.ExpiresAt, &refreshToken.RevokedAt,
	)

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *refreshTokenStore) Revoke(ctx context.Context, id string) error {
	q := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1;`

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrRefreshTokenNotFound
	}

	return nil
}

func (s *refreshTokenStore) RevokeAllByUserID(ctx context.Context, id string) error {
	q := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1;`

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

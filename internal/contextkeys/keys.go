package contextkeys

import (
	"context"
	"errors"
)

type contextKey string

var (
	userIDKey contextKey = "userID"
)

var (
	ErrInvalidCtxVal = errors.New("invalid context value")
)

func SetUserIDCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserIDFromCtx(ctx context.Context) (string, error) {
	userIDCtxVal := ctx.Value(userIDKey)
	userID, ok := userIDCtxVal.(string)
	if !ok || userID == "" {
		return "", ErrInvalidCtxVal
	}
	return userID, nil
}

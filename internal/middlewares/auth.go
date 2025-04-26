package middlewares

import (
	"net/http"
	"strings"

	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/services"
)

type AuthMiddleware interface {
	Authorise(http.Handler) http.Handler
}

type authMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) AuthMiddleware {
	return &authMiddleware{
		authService: authService,
	}
}

func (m *authMiddleware) Authorise(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer ") {
			http.Error(w, "unathorised", http.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer "))
		if accessToken == "" {
			http.Error(w, "unathorised", http.StatusUnauthorized)
			return
		}

		userID, err := m.authService.VerifyAccessToken(r.Context(), accessToken)
		if err != nil {
			http.Error(w, "unathorised", http.StatusUnauthorized)
			return
		}

		ctx := contextkeys.SetUserIDCtx(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

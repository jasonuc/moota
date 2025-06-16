package middlewares

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/services"
)

type AuthMiddleware interface {
	Authorise(http.Handler) http.Handler
	ValidateUserAccess(http.Handler) http.Handler
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
		accessToken, err := r.Cookie("access_token")

		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
			default:
				http.Error(w, "unauthorised", http.StatusUnauthorized)
			}
			return
		}

		if accessToken.Value == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := m.authService.VerifyAccessToken(r.Context(), accessToken.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := contextkeys.SetUserIDCtx(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *authMiddleware) ValidateUserAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDFromCtx, err := contextkeys.GetUserIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		userIDParam := chi.URLParam(r, "userID")
		if userIDFromCtx != userIDParam {
			http.Error(w, "cannot access another user's data", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

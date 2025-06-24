package middlewares

import (
	"net/http"

	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/services"
	"github.com/jasonuc/moota/internal/utils"
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
			utils.UnauthorizedResponse(w)
			return
		}

		if accessToken.Value == "" {
			utils.UnauthorizedResponse(w)
			return
		}

		userID, err := m.authService.VerifyAccessToken(r.Context(), accessToken.Value)
		if err != nil {
			utils.UnauthorizedResponse(w)
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

		userIDParam, err := utils.ReadStringReqParam(r, "userID")
		if err != nil {
			utils.BadRequestResponse(w, err)
			return
		}
		if userIDFromCtx != userIDParam {
			utils.NotPermittedResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

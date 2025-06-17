package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/services"
)

type AuthHandler struct {
	authService  services.AuthService
	validator    *validator.Validate
	cookieDomain string
}

func NewAuthHandler(authService services.AuthService, cookieDomain string) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		validator:    validator.New(),
		cookieDomain: cookieDomain,
	}
}

func (h *AuthHandler) HandleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserRegisterReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	user, tokenPair, err := h.authService.Register(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists):
			badRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidEmail):
			badRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidUsername):
			badRequestResponse(w, err)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Domain:   h.cookieDomain,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Domain:   h.cookieDomain,
		SameSite: http.SameSiteStrictMode,
	})

	//nolint:errcheck
	writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserLoginReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	tokenPair, err := h.authService.Login(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			invalidCredentialsResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		Secure:   true,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		Secure:   true,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	//nolint:errcheck
	writeJSON(w, http.StatusOK, nil, nil)
}

func (h *AuthHandler) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "unauthorised", http.StatusUnauthorized)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	tokenPair, err := h.authService.RefreshTokens(r.Context(), refreshToken.Value)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidRefreshToken):
			invalidCredentialsResponse(w)
		case errors.Is(err, services.ErrTokenExpiredOrRevoked):
			invalidCredentialsResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		MaxAge:   h.authService.GetAccessTokenTTL(),
		Path:     "/",
		Secure:   true,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		MaxAge:   h.authService.GetRefreshTokenTTL(),
		Path:     "/",
		Domain:   h.cookieDomain,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	//nolint:errcheck
	writeJSON(w, http.StatusOK, nil, nil)
}

func (h *AuthHandler) HandleChangeUsername(w http.ResponseWriter, r *http.Request) {
	var payload dto.ChangeUsernameReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	user, err := h.authService.ChangeUserUsername(r.Context(), userID, payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUsername):
			badRequestResponse(w, err)
		case errors.Is(err, models.ErrUserNotFound):
			notFoundResponse(w)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.ChangePasswordReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	user, err := h.authService.ChangeUserPassword(r.Context(), userID, payload)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			notFoundResponse(w)
		case errors.Is(err, services.ErrInvalidCredentials):
			invalidCredentialsResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleChangeEmail(w http.ResponseWriter, r *http.Request) {
	var payload dto.ChangeEmailReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	user, err := h.authService.ChangeUserEmail(r.Context(), userID, payload)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			notFoundResponse(w)
		case errors.Is(err, services.ErrInvalidEmail):
			badRequestResponse(w, err)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	//nolint:errcheck
	writeJSON(w, http.StatusOK, nil, nil)
}

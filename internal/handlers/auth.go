package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/services"
	"github.com/jasonuc/moota/internal/utils"
)

type AuthHandler struct {
	authService        services.AuthService
	validator          *validator.Validate
	cookieDomain       string
	cookieSameSiteMode http.SameSite
}

func NewAuthHandler(authService services.AuthService, cookieDomain string, cookieSameSiteMode int) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		validator:          validator.New(),
		cookieDomain:       cookieDomain,
		cookieSameSiteMode: http.SameSite(cookieSameSiteMode),
	}
}

func (h *AuthHandler) HandleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserRegisterReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	user, tokenPair, err := h.authService.Register(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUsernameTooLong) || errors.Is(err, services.ErrUsernameMustContainOnlyLetters) || errors.Is(err, services.ErrUsernameTaken):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidEmail):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidUsername):
			utils.BadRequestResponse(w, err)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	h.addCookie(w, "access_token", tokenPair.AccessToken, h.authService.GetAccessTokenTTL())
	h.addCookie(w, "refresh_token", tokenPair.RefreshToken, h.authService.GetRefreshTokenTTL())

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserLoginReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	tokenPair, err := h.authService.Login(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.InvalidCredentialsResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	h.addCookie(w, "access_token", tokenPair.AccessToken, h.authService.GetAccessTokenTTL())
	h.addCookie(w, "refresh_token", tokenPair.RefreshToken, h.authService.GetRefreshTokenTTL())

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, nil, nil)
}

func (h *AuthHandler) HandleAccessTokenRefresh(w http.ResponseWriter, r *http.Request) {
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

	tokenPair, err := h.authService.RefreshAccessToken(r.Context(), refreshToken.Value)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidRefreshToken) || errors.Is(err, services.ErrTokenExpiredOrRevoked):
			utils.InvalidCredentialsResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	h.addCookie(w, "access_token", tokenPair.AccessToken, h.authService.GetAccessTokenTTL())
	h.addCookie(w, "refresh_token", tokenPair.RefreshToken, h.authService.GetRefreshTokenTTL())

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, nil, nil)
}

func (h *AuthHandler) HandleChangeUsername(w http.ResponseWriter, r *http.Request) {
	var payload dto.ChangeUsernameReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	userID, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	user, err := h.authService.ChangeUserUsername(r.Context(), userID, payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUsernameTooLong) || errors.Is(err, services.ErrUsernameMustContainOnlyLetters) || errors.Is(err, services.ErrUsernameTaken):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidUsername):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidUsername):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, models.ErrUserNotFound):
			utils.NotFoundResponse(w)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.ChangePasswordReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	userID, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	user, err := h.authService.ChangeUserPassword(r.Context(), userID, payload)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			utils.NotFoundResponse(w)
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.InvalidCredentialsResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleChangeEmail(w http.ResponseWriter, r *http.Request) {
	var payload dto.ChangeEmailReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	userID, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	user, err := h.authService.ChangeUserEmail(r.Context(), userID, payload)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			utils.NotFoundResponse(w)
		case errors.Is(err, services.ErrInvalidEmail):
			utils.BadRequestResponse(w, err)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"user": user}, nil)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	h.deleteCookie(w, "access_token")
	h.deleteCookie(w, "refresh_token")

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, nil, nil)
}

func (h *AuthHandler) addCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: h.cookieSameSiteMode,
	})
}

func (h *AuthHandler) deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Domain:   h.cookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: h.cookieSameSiteMode,
	})
}

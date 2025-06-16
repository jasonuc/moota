package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/services"
)

type UserHandler struct {
	userService services.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := readStringReqParam(r, "userID")
	if err != nil || userID == "" {
		badRequestResponse(w, fmt.Errorf("missing required param userID"))
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			notFoundResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
}

func (h *UserHandler) HandleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	targetUsername, err := readStringReqParam(r, "username")
	if err != nil || targetUsername == "" {
		badRequestResponse(w, fmt.Errorf("missing required param username"))
		return
	}

	userProfile, err := h.userService.GetUserProfile(r.Context(), targetUsername)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			notFoundResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"userProfile": userProfile}, nil)
}

func (h *UserHandler) HandleWhoAmI(w http.ResponseWriter, r *http.Request) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userIDFromCtx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			notFoundResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
}

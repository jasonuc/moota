package handlers

import (
	"errors"
	"math/rand/v2"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/services"
)

type SeedHandler struct {
	seedService services.SeedService
	validator   *validator.Validate
}

func NewSeedHandler(seedService services.SeedService) *SeedHandler {
	return &SeedHandler{
		seedService: seedService,
		validator:   validator.New(),
	}
}

func (h *SeedHandler) HandlePlantSeed(w http.ResponseWriter, r *http.Request) {
	var payload dto.PlantSeedReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	seedID, err := readStringReqParam(r, "seedID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	plant, err := h.seedService.PlantSeed(r.Context(), seedID, payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotPossibleToCreatePlant):
			badRequestResponse(w, err)
		case errors.Is(err, services.ErrNotPossibleToPlantSeed):
			badRequestResponse(w, err)
		case errors.Is(err, services.ErrNotPossibleToPlantSeed):
			badRequestResponse(w, err)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusCreated, envelope{"plant": plant}, nil)
}

func (h *SeedHandler) HandleGetUserSeeds(w http.ResponseWriter, r *http.Request) {
	userIDFromReqParam, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	seedGroups, err := h.seedService.GetAllUserSeeds(r.Context(), userIDFromReqParam)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPermissionsForSeed):
			notPermittedResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"seeds": seedGroups}, nil)
}

func (h *SeedHandler) HandleRequestForNewSeeds(w http.ResponseWriter, r *http.Request) {
	userIDFromReqParam, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	seedGroups, err := h.seedService.GiveUserNewSeeds(r.Context(), userIDFromReqParam, 5+rand.IntN(5)+1) // 5-10 seeds
	if err != nil {
		switch {
		case errors.Is(err, services.ErrSeedRequestInCooldown):
			badRequestResponse(w, err)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"seeds": seedGroups}, nil)
}

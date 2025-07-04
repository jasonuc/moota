package handlers

import (
	"errors"
	"math/rand/v2"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/events"
	"github.com/jasonuc/moota/internal/services"
	"github.com/jasonuc/moota/internal/utils"
)

type SeedHandler struct {
	eventBus    *cqrs.EventBus
	seedService services.SeedService
	validator   *validator.Validate
}

func NewSeedHandler(seedService services.SeedService, eventBus *cqrs.EventBus) *SeedHandler {
	return &SeedHandler{
		seedService: seedService,
		validator:   validator.New(),
		eventBus:    eventBus,
	}
}

func (h *SeedHandler) HandlePlantSeed(w http.ResponseWriter, r *http.Request) {
	var payload dto.PlantSeedReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	seedID, err := utils.ReadStringReqParam(r, "seedID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	plant, err := h.seedService.PlantSeed(r.Context(), seedID, payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotPossibleToCreatePlant):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrNotPossibleToPlantSeed):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrNotPossibleToPlantSeed):
			utils.BadRequestResponse(w, err)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}
	err = h.eventBus.Publish(r.Context(), events.StatUpdated{})
	if err != nil {
		utils.ServerErrorResponse(w, err)
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"plant": plant}, nil)
}

func (h *SeedHandler) HandleGetUserSeeds(w http.ResponseWriter, r *http.Request) {
	userIDFromReqParam, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	seedGroups, err := h.seedService.GetUserSeeds(r.Context(), userIDFromReqParam)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPermissionsForSeed):
			utils.NotPermittedResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"seeds": seedGroups}, nil)
}

func (h *SeedHandler) HandleRequestForNewSeeds(w http.ResponseWriter, r *http.Request) {
	userIDFromReqParam, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	seedGroups, err := h.seedService.GiveUserNewSeeds(r.Context(), userIDFromReqParam, 5+rand.IntN(5)+1) // 5-10 seeds
	if err != nil {
		var errSeedRequestInCooldown *services.ErrSeedRequestInCooldown
		switch {
		case errors.As(err, &errSeedRequestInCooldown):
			utils.ErrorResponse(w, http.StatusForbidden, errSeedRequestInCooldown)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}
	err = h.eventBus.Publish(r.Context(), events.SeedGenerated{})
	if err != nil {
		utils.ServerErrorResponse(w, err)
		return
	}
	//nolint:errcheck
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"seeds": seedGroups}, nil)
}

func (h *SeedHandler) HandleCheckWhenUserCanRequestSeed(w http.ResponseWriter, r *http.Request) {
	userIDFromReqParam, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	timeUntilUserCanReqSeeds, err := h.seedService.CheckWhenUserCanRequestSeed(r.Context(), userIDFromReqParam)
	if err != nil {
		utils.ServerErrorResponse(w, err)
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"timeAvailable": timeUntilUserCanReqSeeds, "availableNow": timeUntilUserCanReqSeeds == nil}, nil)
}

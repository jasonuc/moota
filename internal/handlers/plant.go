package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/services"
	"github.com/jasonuc/moota/internal/store"
)

type PlantHandler struct {
	plantService services.PlantService
	validator    *validator.Validate
}

func NewPlantService(plantService services.PlantService) *PlantHandler {
	return &PlantHandler{
		plantService: plantService,
		validator:    validator.New(),
	}
}

func (h *PlantHandler) HandleGetAllUserPlants(w http.ResponseWriter, r *http.Request) {
	lon, err := readFloatQueryParam(r, "lon")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	lat, err := readFloatQueryParam(r, "lat")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	userCoords := &models.Coordinates{Lat: lat, Lon: lon}

	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	IncludeDeceased := readBoolQueryParam(r, "includeDeceased")

	plants, err := h.plantService.GetAllUserPlants(r.Context(), userID, userCoords, &store.GetPlantsOpts{IncludeDeceased: IncludeDeceased})
	if err != nil {
		serverErrorResponse(w, err)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"plants": plants}, nil)
}

func (h *PlantHandler) HandleGetPlant(w http.ResponseWriter, r *http.Request) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(r.Context())
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	plant, err := h.plantService.GetPlant(r.Context(), plantID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPlantNotFound):
			notFoundResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	if plant.OwnerID != userIDFromCtx {
		plant.C = models.Coordinates{}
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleActionOnPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	var payload dto.ActionOnPlantReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	plant, err := h.plantService.ActionOnPlant(r.Context(), plantID, payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorisedPlantAction):
			notPermittedResponse(w)
		case errors.Is(err, models.ErrPlantNotFound):
			notFoundResponse(w)
		case errors.Is(err, services.ErrOutsidePlantInteractionRadius):
			badRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidPlantAction):
			badRequestResponse(w, err)
		case errors.Is(err, models.ErrPlantInCooldown):
			badRequestResponse(w, err)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleGetAllUserDeceasedPlants(w http.ResponseWriter, r *http.Request) {
	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	deceasedPlants, err := h.plantService.GetAllUserDeceasedPlants(r.Context(), userID)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"plants": deceasedPlants}, nil)
}

func (h *PlantHandler) HandleKillPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.plantService.KillPlant(r.Context(), plantID); err != nil {
		switch {
		case errors.Is(err, models.ErrPlantNotFound):
			notFoundResponse(w)
		case errors.Is(err, services.ErrPlantAlreadyDead):
			badRequestResponse(w, err)
		case errors.Is(err, services.ErrUnauthorisedPlantAction):
			notPermittedResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, nil, nil)
}

func (h *PlantHandler) HandleChangePlantNickname(w http.ResponseWriter, r *http.Request) {
	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		badRequestResponse(w, err)
		return
	}

	var payload dto.ChangePlantNicknameReq
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		failedValidationResponse(w, err)
		return
	}

	plant, err := h.plantService.ChangePlantNickname(r.Context(), plantID, payload.NewNickname)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPlantNotFound):
			notFoundResponse(w)
		case errors.Is(err, services.ErrUnauthorisedPlantAction):
			notPermittedResponse(w)
		default:
			serverErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"plant": plant}, nil)
}

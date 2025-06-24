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
	"github.com/jasonuc/moota/internal/utils"
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
	lon, err := utils.ReadFloatQueryParam(r, "lon")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	lat, err := utils.ReadFloatQueryParam(r, "lat")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	userCoords := &models.Coordinates{Lat: lat, Lon: lon}

	userID, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	IncludeDeceased := utils.ReadBoolQueryParam(r, "includeDeceased")

	plants, err := h.plantService.GetAllUserPlants(r.Context(), userID, userCoords, &store.GetPlantsOpts{IncludeDeceased: IncludeDeceased})
	if err != nil {
		utils.ServerErrorResponse(w, err)
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"plants": plants}, nil)
}

func (h *PlantHandler) HandleGetPlant(w http.ResponseWriter, r *http.Request) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(r.Context())
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	plantID, err := utils.ReadStringReqParam(r, "plantID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	plant, err := h.plantService.GetPlant(r.Context(), plantID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPlantNotFound):
			utils.NotFoundResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	if plant.OwnerID != userIDFromCtx {
		plant.C = models.Coordinates{}
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleActionOnPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := utils.ReadStringReqParam(r, "plantID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	var payload dto.ActionOnPlantReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	plant, err := h.plantService.ActionOnPlant(r.Context(), plantID, payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorisedPlantAction):
			utils.NotPermittedResponse(w)
		case errors.Is(err, models.ErrPlantNotFound):
			utils.NotFoundResponse(w)
		case errors.Is(err, services.ErrOutsidePlantInteractionRadius):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrInvalidPlantAction):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, models.ErrPlantInCooldown):
			utils.BadRequestResponse(w, err)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleGetAllUserDeceasedPlants(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadStringReqParam(r, "userID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	deceasedPlants, err := h.plantService.GetAllUserDeceasedPlants(r.Context(), userID)
	if err != nil {
		utils.ServerErrorResponse(w, err)
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"plants": deceasedPlants}, nil)
}

func (h *PlantHandler) HandleKillPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := utils.ReadStringReqParam(r, "plantID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.plantService.KillPlant(r.Context(), plantID); err != nil {
		switch {
		case errors.Is(err, models.ErrPlantNotFound):
			utils.NotFoundResponse(w)
		case errors.Is(err, services.ErrPlantAlreadyDead):
			utils.BadRequestResponse(w, err)
		case errors.Is(err, services.ErrUnauthorisedPlantAction):
			utils.NotPermittedResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusAccepted, nil, nil)
}

func (h *PlantHandler) HandleChangePlantNickname(w http.ResponseWriter, r *http.Request) {
	plantID, err := utils.ReadStringReqParam(r, "plantID")
	if err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	var payload dto.ChangePlantNicknameReq
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, err)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		utils.FailedValidationResponse(w, err)
		return
	}

	plant, err := h.plantService.ChangePlantNickname(r.Context(), plantID, payload.NewNickname)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPlantNotFound):
			utils.NotFoundResponse(w)
		case errors.Is(err, services.ErrUnauthorisedPlantAction):
			utils.NotPermittedResponse(w)
		default:
			utils.ServerErrorResponse(w, err)
		}
		return
	}

	//nolint:errcheck
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"plant": plant}, nil)
}

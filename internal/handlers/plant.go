package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/services"
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
	var payload dto.GetAllUserPlantsReq
	if err := readJSON(w, r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	plants, err := h.plantService.GetAllUserPlants(r.Context(), userID, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusOK, envelope{"plants": plants}, nil)
}

func (h *PlantHandler) HandleGetPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.GetPlant(r.Context(), plantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleActionOnPlant(w http.ResponseWriter, r *http.Request) {
	var payload dto.ActionOnPlantReq
	if err := readJSON(w, r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.ActionOnPlant(r.Context(), plantID, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleActivatePlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plant, err := h.plantService.ActivatePlant(r.Context(), plantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"plant": plant}, nil)
}

func (h *PlantHandler) HandleGetAllUserDeceasedPlants(w http.ResponseWriter, r *http.Request) {
	userID, err := readStringReqParam(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deceasedPlants, err := h.plantService.GetAllUserDeceasedPlants(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, envelope{"plants": deceasedPlants}, nil)
}

func (h *PlantHandler) HandleKillPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := readStringReqParam(r, "plantID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.plantService.KillPlant(r.Context(), plantID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//nolint:errcheck
	writeJSON(w, http.StatusAccepted, nil, nil)
}

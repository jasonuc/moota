package handlers

import (
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plant, err := h.seedService.PlantSeed(r.Context(), payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, envelope{"plant": plant}, nil)
}

func (h *SeedHandler) HandleGetUserSeeds(w http.ResponseWriter, r *http.Request) {
	userIDFromReqParam, err := readStringReqParam(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seedGroups, err := h.seedService.GetAllUserSeeds(r.Context(), userIDFromReqParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, envelope{"seeds": seedGroups}, nil)
}

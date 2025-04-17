package services

import (
	"errors"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type PlantService struct {
	store *store.Store
}

func NewPlantService(store *store.Store) *PlantService {
	return &PlantService{
		store: store,
	}
}

var (
	ErrNotPossibleToCreatePlant = errors.New("not possible to create plant here")
)

func (s *PlantService) CreatePlant(tx *store.Store, soil *models.Soil, seed *models.Seed, centre models.Coordinates) (*models.Plant, error) {
	plantCircleMeta := models.NewCircleMeta(soil.Centre(), models.PlantInteractionRadius)
	nearbyPlants, err := s.store.Plant.GetAllInSoilAndInProximity(soil.ID, centre, models.PlantInteractionRadius+1)
	if err != nil {
		return nil, err
	}
	if !s.IsPlantValidForSoil(plantCircleMeta, nearbyPlants) {
		return nil, ErrNotPossibleToCreatePlant
	}

	plant, err := models.NewPlant(seed, soil, centre)
	if err != nil {
		return nil, err
	}

	err = tx.Plant.Insert(plant)
	if err != nil {
		return nil, err
	}

	return plant, nil
}

func (s *PlantService) IsPlantValidForSoil(plantCircleMeta models.CircleMeta, nearbyPlants []*models.Plant) bool {
	plantsOverlapMap := make(map[bool]struct{})
	for _, nearbyPlant := range nearbyPlants {
		plantsOverlapMap[plantCircleMeta.OverlapsWith(nearbyPlant)] = struct{}{}
	}

	_, ok := plantsOverlapMap[true]
	return !ok
}

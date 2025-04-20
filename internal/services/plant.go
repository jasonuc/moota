package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type PlantService struct {
	store *store.Store
}

type ActionOnPlantReqDto struct {
	Action    int
	Longitude float64
	Latitude  float64
	Time      time.Time
	UserID    string
	PlantID   string
}

func NewPlantService(store *store.Store) *PlantService {
	return &PlantService{
		store: store,
	}
}

func (s *PlantService) WithStore(store *store.Store) *PlantService {
	copy := *s
	copy.store = store
	return &copy
}

var (
	ErrNotPossibleToCreatePlant      = errors.New("not possible to create plant here")
	ErrOutsidePlantInteractionRadius = errors.New("user is not within plant interaction radius")
	ErrInvalidPlantAction            = errors.New("invalid plant action")
)

func (s *PlantService) ConfirmPlantCreation(plantID string) (*models.Plant, error) {
	plant, err := s.store.Plant.Get(plantID)
	if err != nil {
		return nil, err
	}
	if plant.Activated {
		return nil, fmt.Errorf("plant already activated")
	}

	if err := s.store.Plant.ActivatePlant(plant.ID); err != nil {
		return nil, err
	}
	return s.store.Plant.Get(plant.ID)
}

func (s *PlantService) GetAllUserPlants(ownerID string) ([]*models.Plant, error) {
	plants, err := s.store.Plant.GetAllByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	return plants, nil
}

func (s *PlantService) GetPlant(ownerID, plantID string) (*models.Plant, error) {
	plant, err := s.store.Plant.Get(plantID)
	if err != nil {
		return nil, err
	}
	if plant.OwnerID != ownerID {
		return nil, fmt.Errorf("access denied: you do not own this plant")
	}
	return plant, nil
}

func (s *PlantService) CreatePlant(soil *models.Soil, seed *models.Seed, centre models.Coordinates) (*models.Plant, error) {
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

	err = s.store.Plant.Insert(plant)
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

func (s *PlantService) ActionOnPlant(dto ActionOnPlantReqDto) (*models.Plant, error) {
	if !models.ValidPlantAction(dto.Action) {
		return nil, ErrInvalidPlantAction
	}

	plant, err := s.GetPlant(dto.UserID, dto.PlantID)
	if err != nil {
		return nil, err
	}

	userCoords := models.Coordinates{Lng: dto.Longitude, Lat: dto.Latitude}
	if !plant.ContainsPoint(userCoords) {
		return nil, ErrOutsidePlantInteractionRadius
	}

	_, err = plant.Action(models.PlantAction(dto.Action), dto.Time)
	if err != nil {
		return nil, err
	}

	err = s.store.Plant.Update(plant)
	if err != nil {
		return nil, err
	}

	return s.GetPlant(plant.OwnerID, dto.PlantID)
}

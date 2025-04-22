package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type PlantService interface {
	GetAllUserPlants(string, models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error)
	Get4ClosestPlants(string, models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error)
	GetPlant(string, string) (*models.Plant, error)
	CreatePlant(*models.Soil, *models.Seed, models.Coordinates) (*models.Plant, error)
	KillPlant(string) error
	WithStore(*store.Store) PlantService
}

type plantService struct {
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

func NewPlantService(store *store.Store) PlantService {
	return &plantService{
		store: store,
	}
}

func (s *plantService) WithStore(store *store.Store) PlantService {
	copy := *s
	copy.store = store
	return &copy
}

var (
	ErrNotPossibleToCreatePlant      = errors.New("not possible to create plant here")
	ErrOutsidePlantInteractionRadius = errors.New("user is not within plant interaction radius")
	ErrInvalidPlantAction            = errors.New("invalid plant action")
	ErrPlantAlreadyActivated         = errors.New("plant already activated")
	ErrPlantAreadyDead               = errors.New("plant already dead")
)

func (s *plantService) GetAllUserPlants(userID string, point models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plants, err := tx.Plant.GetAllByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	refreshPlantsData(tx, plants, now)

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	plantsWithDistanceM := make([]*models.PlantWithDistanceMFromUser, 0)
	for _, p := range plants {
		plantsWithDistanceM = append(plantsWithDistanceM, &models.PlantWithDistanceMFromUser{
			Plant:     *p,
			DistanceM: p.Centre().DistanceM(point),
		})
	}

	return plantsWithDistanceM, nil
}

func (s *plantService) Get4ClosestPlants(userID string, point models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error) {
	plants, err := s.GetAllUserPlants(userID, point)
	if err != nil {
		return nil, err
	}
	return plants[:4], err
}

func (s *plantService) GetPlant(userID, plantID string) (*models.Plant, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(plantID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	refreshPlantData(tx, plant, now)

	if plant.OwnerID != userID {
		return nil, fmt.Errorf("access denied: you do not own this plant")
	}
	return plant, nil
}

func (s *plantService) ConfirmPlantCreation(plantID string) (*models.Plant, error) {
	plant, err := s.store.Plant.Get(plantID)
	if err != nil {
		return nil, err
	}
	if plant.Activated {
		return nil, ErrPlantAlreadyActivated
	}

	if err := s.store.Plant.ActivatePlant(plant.ID); err != nil {
		return nil, err
	}
	return s.store.Plant.Get(plant.ID)
}

func (s *plantService) CreatePlant(soil *models.Soil, seed *models.Seed, centre models.Coordinates) (*models.Plant, error) {
	plantCircleMeta := models.NewCircleMeta(soil.Centre(), models.PlantInteractionRadius)
	nearbyPlants, err := s.store.Plant.GetAllInSoilAndInProximity(soil.ID, centre, models.PlantInteractionRadius+1)
	if err != nil {
		return nil, err
	}
	if !s.isPlantValidForSoil(plantCircleMeta, nearbyPlants) {
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

func (s *plantService) ActionOnPlant(dto ActionOnPlantReqDto) (*models.Plant, error) {
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

func (s *plantService) KillPlant(id string) error {
	transaction, err := s.store.Begin()
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(id)
	if err != nil {
		return err
	}

	if plant.Dead {
		return ErrPlantAreadyDead
	}

	plant.Dead = true
	if err := tx.Plant.Update(plant); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func refreshPlantsData(tx *store.Store, plants []*models.Plant, t time.Time) {
	for _, plant := range plants {
		refreshPlantData(tx, plant, t)
	}
}

func refreshPlantData(tx *store.Store, plant *models.Plant, t time.Time) {
	plant.Refresh(t)
	tx.Plant.Update(plant)
}

func (s *plantService) isPlantValidForSoil(plantCircleMeta models.CircleMeta, nearbyPlants []*models.Plant) bool {
	plantsOverlapMap := make(map[bool]struct{})
	for _, nearbyPlant := range nearbyPlants {
		plantsOverlapMap[plantCircleMeta.OverlapsWith(nearbyPlant)] = struct{}{}
	}

	_, ok := plantsOverlapMap[true]
	return !ok
}

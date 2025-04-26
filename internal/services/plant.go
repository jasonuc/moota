package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type PlantService interface {
	GetAllUserPlants(context.Context, string, models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error)
	Get4ClosestPlants(context.Context, string, models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error)
	GetPlant(context.Context, string, string) (*models.Plant, error)
	CreatePlant(context.Context, *models.Soil, *models.Seed, models.Coordinates) (*models.Plant, error)
	KillPlant(context.Context, string) error
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

func (s *plantService) GetAllUserPlants(ctx context.Context, userID string, point models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plants, err := tx.Plant.GetByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	err = refreshPlantsData(ctx, tx, plants, now)
	if err != nil {
		return nil, err
	}

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

func (s *plantService) Get4ClosestPlants(ctx context.Context, userID string, point models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error) {
	plants, err := s.GetAllUserPlants(ctx, userID, point)
	if err != nil {
		return nil, err
	}
	return plants[:4], err
}

func (s *plantService) GetPlant(ctx context.Context, userID, plantID string) (*models.Plant, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(ctx, plantID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	err = refreshPlantData(ctx, tx, plant, now)
	if err != nil {
		return nil, err
	}

	if plant.OwnerID != userID {
		return nil, fmt.Errorf("access denied: you do not own this plant")
	}
	return plant, nil
}

func (s *plantService) ConfirmPlantCreation(ctx context.Context, plantID string) (*models.Plant, error) {
	plant, err := s.store.Plant.Get(ctx, plantID)
	if err != nil {
		return nil, err
	}
	if plant.Activated {
		return nil, ErrPlantAlreadyActivated
	}

	if err := s.store.Plant.ActivatePlant(ctx, plant.ID); err != nil {
		return nil, err
	}
	return s.store.Plant.Get(ctx, plant.ID)
}

func (s *plantService) CreatePlant(ctx context.Context, soil *models.Soil, seed *models.Seed, centre models.Coordinates) (*models.Plant, error) {
	plantCircleMeta := models.NewCircleMeta(soil.Centre(), models.PlantInteractionRadius)
	nearbyPlants, err := s.store.Plant.GetBySoilIDAndProximity(ctx, soil.ID, centre, models.PlantInteractionRadius+1)
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

	err = s.store.Plant.Insert(ctx, plant)
	if err != nil {
		return nil, err
	}

	return plant, nil
}

func (s *plantService) ActionOnPlant(ctx context.Context, dto ActionOnPlantReqDto) (*models.Plant, error) {
	if !models.ValidPlantAction(dto.Action) {
		return nil, ErrInvalidPlantAction
	}

	plant, err := s.GetPlant(ctx, dto.UserID, dto.PlantID)
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

	err = s.store.Plant.Update(ctx, plant)
	if err != nil {
		return nil, err
	}

	return s.GetPlant(ctx, plant.OwnerID, dto.PlantID)
}

func (s *plantService) KillPlant(ctx context.Context, id string) error {
	transaction, err := s.store.Begin()
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(ctx, id)
	if err != nil {
		return err
	}

	if plant.Dead {
		return ErrPlantAreadyDead
	}

	plant.Dead = true
	if err := tx.Plant.Update(ctx, plant); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func refreshPlantsData(ctx context.Context, tx *store.Store, plants []*models.Plant, t time.Time) error {
	for _, plant := range plants {
		err := refreshPlantData(ctx, tx, plant, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshPlantData(ctx context.Context, tx *store.Store, plant *models.Plant, t time.Time) error {
	plant.Refresh(t)
	return tx.Plant.Update(ctx, plant)
}

func (s *plantService) isPlantValidForSoil(plantCircleMeta models.CircleMeta, nearbyPlants []*models.Plant) bool {
	plantsOverlapMap := make(map[bool]struct{})
	for _, nearbyPlant := range nearbyPlants {
		plantsOverlapMap[plantCircleMeta.OverlapsWith(nearbyPlant)] = struct{}{}
	}

	_, ok := plantsOverlapMap[true]
	return !ok
}

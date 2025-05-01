package services

import (
	"context"
	"errors"
	"time"

	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type PlantService interface {
	GetAllUserPlants(context.Context, string, dto.GetAllUserPlantsReq) ([]*models.PlantWithDistanceMFromUser, error)
	ActionOnPlant(context.Context, string, dto.ActionOnPlantReq) (*models.Plant, error)
	GetPlant(context.Context, string) (*models.Plant, error)
	CreatePlant(context.Context, *models.Soil, *models.Seed, models.Coordinates) (*models.Plant, error)
	GetAllUserDeceasedPlants(context.Context, string) ([]*models.Plant, error)
	ActivatePlant(ctx context.Context, plantID string) (*models.Plant, error)
	KillPlant(context.Context, string) error
	WithStore(*store.Store) PlantService
}

type plantService struct {
	store *store.Store
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
	ErrUnauthorisedPlantAction       = errors.New("unauthorised plant action")
	ErrPlantAlreadyActivated         = errors.New("plant already activated")
	ErrPlantAreadyDead               = errors.New("plant already dead")
	ErrPlantNotActivated             = errors.New("plant not activated")
)

func (s *plantService) GetAllUserPlants(ctx context.Context, userID string, dto dto.GetAllUserPlantsReq) ([]*models.PlantWithDistanceMFromUser, error) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	coords := models.Coordinates{Lat: *dto.Latitude, Lng: *dto.Longitude}

	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	if userID != userIDFromCtx {
		return nil, ErrUnauthorisedPlantAction
	}

	tx := s.store.WithTx(transaction)

	plants, err := tx.Plant.GetByOwnerID(ctx, userID, false)
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
			DistanceM: p.Centre().DistanceM(coords),
		})
	}

	return plantsWithDistanceM, nil
}

func (s *plantService) GetPlant(ctx context.Context, plantID string) (*models.Plant, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(ctx, plantID, false)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	err = refreshPlantData(ctx, tx, plant, now)
	if err != nil {
		return nil, err
	}

	return plant, nil
}

func (s *plantService) ActivatePlant(ctx context.Context, plantID string) (*models.Plant, error) {
	plant, err := s.store.Plant.Get(ctx, plantID, false)
	if err != nil {
		return nil, err
	}
	if plant.Activated {
		return nil, ErrPlantAlreadyActivated
	}

	if err := s.store.Plant.ActivatePlant(ctx, plant.ID); err != nil {
		return nil, err
	}
	return s.store.Plant.Get(ctx, plant.ID, false)
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

func (s *plantService) ActionOnPlant(ctx context.Context, plantID string, dto dto.ActionOnPlantReq) (*models.Plant, error) {
	userID, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !models.ValidPlantAction(dto.Action) {
		return nil, ErrInvalidPlantAction
	}

	plant, err := s.GetPlant(ctx, plantID)
	if err != nil {
		return nil, err
	}

	if plant.OwnerID != userID {
		return nil, ErrUnauthorisedPlantAction
	}

	if !plant.Activated {
		return nil, ErrPlantNotActivated
	}

	userCoords := models.Coordinates{Lng: *dto.Longitude, Lat: *dto.Latitude}
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

	return s.GetPlant(ctx, plantID)
}

func (s *plantService) GetAllUserDeceasedPlants(ctx context.Context, userID string) ([]*models.Plant, error) {
	allUserPlants, err := s.store.Plant.GetByOwnerID(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	deceasedPlants := make([]*models.Plant, 0)
	for _, plant := range allUserPlants {
		if plant.Dead {
			deceasedPlants = append(deceasedPlants, plant)
		}
	}

	return deceasedPlants, nil
}

func (s *plantService) KillPlant(ctx context.Context, id string) error {
	transaction, err := s.store.Begin()
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(ctx, id, false)
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

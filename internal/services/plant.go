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
	GetAllUserPlants(context.Context, string, *models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error)
	ActionOnPlant(context.Context, string, dto.ActionOnPlantReq) (*models.Plant, error)
	GetPlant(context.Context, string) (*models.Plant, error)
	CreatePlant(context.Context, *models.Soil, *models.Seed, models.Coordinates) (*models.Plant, error)
	GetAllUserDeceasedPlants(context.Context, string) ([]*models.Plant, error)
	ChangePlantNickname(context.Context, string, string) (*models.Plant, error)
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
	ErrNotPossibleToCreatePlant                  = errors.New("not possible to create plant here")
	ErrNotPossibleToActivatePlant                = errors.New("not possible to activate plant")
	ErrNotPossibleToActivatePlantButSeedRefunded = errors.New("not possible to activate plant but seed refunded")
	ErrOutsidePlantInteractionRadius             = errors.New("user is not within plant interaction radius")
	ErrInvalidPlantAction                        = errors.New("invalid plant action")
	ErrUnauthorisedPlantAction                   = errors.New("unauthorised plant action")
	ErrPlantAlreadyActivated                     = errors.New("plant already activated")
	ErrPlantAlreadyDead                          = errors.New("plant already dead")
	ErrPlantNotActivated                         = errors.New("plant not activated")
)

func (s *plantService) GetAllUserPlants(ctx context.Context, userID string, dto *models.Coordinates) ([]*models.PlantWithDistanceMFromUser, error) {
	coords := models.Coordinates{Lat: dto.Lat, Lon: dto.Lon}

	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

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
		return nil, store.ErrTransactionCouldNotStart
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
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()
	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(ctx, plantID, false)
	if err != nil {
		return nil, err
	}

	plantsWithinInteractionRadius, err := tx.Plant.GetBySoilIDAndProximity(ctx, plant.Soil.ID, plant.Centre(), models.PlantInteractionRadius)
	if err != nil {
		return nil, err
	}

	if len(plantsWithinInteractionRadius) != 0 {
		if plant.ID != plantsWithinInteractionRadius[0].ID {
			newSeed := models.NewSeedWithMeta(plant.OwnerID, plant.SeedMeta)
			if err := tx.Seed.Insert(ctx, newSeed); err != nil {
				return nil, err
			}

			if err := tx.Plant.Delete(ctx, plant.ID); err != nil {
				return nil, err
			}

			if err := transaction.Commit(); err != nil {
				return nil, err
			}

			return nil, ErrNotPossibleToActivatePlantButSeedRefunded
		}
		return nil, ErrNotPossibleToActivatePlant
	}

	if plant.Activated {
		return nil, ErrPlantAlreadyActivated
	}

	if err := tx.Plant.ActivatePlant(ctx, plant.ID); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return s.store.Plant.Get(ctx, plant.ID, false)
}

func (s *plantService) CreatePlant(ctx context.Context, soil *models.Soil, seed *models.Seed, centre models.Coordinates) (*models.Plant, error) {
	plantCircleMeta := models.NewCircleMeta(soil.Centre(), models.PlantInteractionRadius)
	nearbyPlants, err := s.store.Plant.GetBySoilIDAndProximity(ctx, soil.ID, centre, models.PlantInteractionRadius+0.1)
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

	userCoords := models.Coordinates{Lon: *dto.Longitude, Lat: *dto.Latitude}
	if !plant.ContainsPoint(userCoords) {
		return nil, ErrOutsidePlantInteractionRadius
	}

	_, err = plant.Action(models.PlantAction(dto.Action), time.Now())
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

func (s *plantService) ChangePlantNickname(ctx context.Context, plantID string, newNickname string) (*models.Plant, error) {
	userID, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	plant, err := tx.Plant.Get(ctx, plantID, false)
	if err != nil {
		return nil, err
	}

	if err := refreshPlantData(ctx, tx, plant, time.Now()); err != nil {
		return nil, err
	}

	if plant.OwnerID != userID {
		return nil, ErrUnauthorisedPlantAction
	}

	plant.Nickname = newNickname
	if err := tx.Plant.Update(ctx, plant); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return plant, nil
}

func (s *plantService) KillPlant(ctx context.Context, id string) error {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return err
	}

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

	if plant.OwnerID != userIDFromCtx {
		return ErrUnauthorisedPlantAction
	}

	if plant.Dead {
		return ErrPlantAlreadyDead
	}

	plant.Die(time.Now())
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

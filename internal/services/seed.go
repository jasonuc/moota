package services

import (
	"context"
	"errors"
	"maps"
	"slices"
	"time"

	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SeedService interface {
	GetUserSeeds(context.Context, string) ([]*models.SeedGroup, error)
	GetSeed(context.Context, string, string) (*models.Seed, error)
	GiveUserNewSeeds(context.Context, string, int) ([]*models.SeedGroup, error)
	PlantSeed(context.Context, string, dto.PlantSeedReq) (*models.Plant, error)
	CheckWhenUserCanRequestSeed(ctx context.Context, userID string) (*time.Time, error)
	WithStore(*store.Store) SeedService
}

var (
	ErrUnauthorizedSeedPlanting  = errors.New("not authorised to plant this seed")
	ErrNotPossibleToPlantSeed    = errors.New("not possible to plant seed")
	ErrInvalidPermissionsForSeed = errors.New("invalid permissions to retreive seed")
)

type ErrSeedRequestInCooldown struct {
	Message       string    `json:"message"`
	TimeAvailable time.Time `json:"timeAvailable"`
}

func (e *ErrSeedRequestInCooldown) Error() string {
	return e.Message
}

var (
	SeedRequestCooldownDuration = 7 * 24 * time.Hour
)

type seedService struct {
	soilService  SoilService
	plantService PlantService
	store        *store.Store
}

func NewSeedService(store *store.Store, soilService SoilService, plantService PlantService) SeedService {
	return &seedService{
		store:        store,
		soilService:  soilService,
		plantService: plantService,
	}
}

func (s *seedService) WithStore(store *store.Store) SeedService {
	copy := *s
	copy.store = store
	return &copy
}

func (s *seedService) GetUserSeeds(ctx context.Context, userID string) ([]*models.SeedGroup, error) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if userID != userIDFromCtx {
		return nil, ErrInvalidPermissionsForSeed
	}

	seeds, err := s.store.Seed.GetByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	seedGroupsMap := make(map[string]*models.SeedGroup)
	for _, seed := range seeds {
		sg, ok := seedGroupsMap[seed.BotanicalName]
		if ok {
			sg.Seeds = append(sg.Seeds, seed)
			sg.Count++
			continue
		} else {
			seedGroupsMap[seed.BotanicalName] = &models.SeedGroup{
				BotanicalName: seed.BotanicalName,
				Count:         1,
				Seeds:         []*models.Seed{seed},
			}
		}
	}

	seedGroups := slices.Collect(maps.Values(seedGroupsMap))
	return seedGroups, nil
}

func (s *seedService) GetSeed(ctx context.Context, userID, seedID string) (*models.Seed, error) {
	seed, err := s.store.Seed.Get(ctx, seedID)
	if err != nil {
		return nil, err
	}
	if seed.OwnerID == userID {
		return nil, ErrInvalidPermissionsForSeed
	}
	return seed, nil
}

func (s *seedService) PlantSeed(ctx context.Context, seedID string, dto dto.PlantSeedReq) (*models.Plant, error) {
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
	soilServiceWithTx := s.soilService.WithStore(tx)
	plantServiceWithTx := s.plantService.WithStore(tx)

	seed, err := tx.Seed.Get(ctx, seedID)
	if err != nil {
		return nil, err
	}

	if seed.OwnerID != userID {
		return nil, ErrUnauthorizedSeedPlanting
	}

	targetCentre := models.Coordinates{Lat: *dto.Latitude, Lon: *dto.Longitude}
	plantCircleMeta := models.NewCircleMeta(targetCentre, models.PlantInteractionRadius)

	nearbySoils, err := tx.Soil.GetAllInProximity(ctx, targetCentre, models.SoilRadiusMLarge)
	if err != nil {
		return nil, err
	}

	if len(nearbySoils) == 0 {
		soil, err := soilServiceWithTx.CreateSoil(ctx, targetCentre, nearbySoils)
		if err != nil {
			return nil, err
		}
		nearbySoils = append(nearbySoils, soil)
	}

	var targetSoil *models.Soil = nil
	for _, soil := range nearbySoils {
		if soil.ContainsFullCircle(plantCircleMeta) {
			targetSoil = soil
		}
	}

	if targetSoil == nil {
		return nil, ErrNotPossibleToPlantSeed
	}

	plant, err := plantServiceWithTx.CreatePlant(ctx, targetSoil, seed, targetCentre)
	if err != nil {
		return nil, err
	}

	if err := tx.Seed.MarkAsPlanted(ctx, seed.ID); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return plant, nil
}

func (s *seedService) CheckWhenUserCanRequestSeed(ctx context.Context, userID string) (*time.Time, error) {
	lastFulfilledSeedRequest, err := s.store.Seed.GetLastFulfilledSeedRequestTimeByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if lastFulfilledSeedRequest.IsZero() {
		return nil, nil
	}

	timeAvailable := lastFulfilledSeedRequest.Add(SeedRequestCooldownDuration)
	if time.Since(lastFulfilledSeedRequest) < SeedRequestCooldownDuration {
		return &timeAvailable, nil
	}

	return nil, nil
}

func (s *seedService) GiveUserNewSeeds(ctx context.Context, userID string, count int) ([]*models.SeedGroup, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)

	lastFulfilledSeedRequest, err := tx.Seed.GetLastFulfilledSeedRequestTimeByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !lastFulfilledSeedRequest.IsZero() {
		if time.Since(lastFulfilledSeedRequest) < SeedRequestCooldownDuration {
			if err := recordFailedSeedRequest(ctx, transaction, tx, userID, count); err != nil {
				return nil, err
			}
			timeAvailable := lastFulfilledSeedRequest.Add(SeedRequestCooldownDuration)
			return nil, &ErrSeedRequestInCooldown{TimeAvailable: timeAvailable, Message: "seed request in cooldown"}
		}
	}

	for range count {
		newSeed := models.NewSeed(userID)
		if err := tx.Seed.Insert(ctx, newSeed); err != nil {
			return nil, err
		}
	}

	if err := tx.Seed.InsertSeedRequest(ctx, userID, time.Now(), true, count); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return s.GetUserSeeds(ctx, userID)
}

func recordFailedSeedRequest(ctx context.Context, transaction *store.Transaction, txStore *store.Store, userID string, count int) error {
	if err := txStore.Seed.InsertSeedRequest(ctx, userID, time.Now(), false, count); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
		return err
	}
	return nil
}

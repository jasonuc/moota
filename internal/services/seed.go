package services

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/jasonuc/moota/internal/contextkeys"
	"github.com/jasonuc/moota/internal/dto"
	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SeedService interface {
	GetAllUserSeeds(context.Context, string) ([]*SeedGroup, error)
	GetSeed(context.Context, string, string) (*models.Seed, error)
	PlantSeed(context.Context, string, dto.PlantSeedReq) (*models.Plant, error)
	WithStore(*store.Store) SeedService
}

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

type SeedGroup struct {
	BotanicalName string         `json:"botanicalName"`
	Count         int            `json:"count"`
	Seeds         []*models.Seed `json:"seeds"`
}

var (
	ErrUnauthorizedSeedPlanting  = errors.New("not authorised to plant this seed")
	ErrNotPossibleToPlantSeed    = errors.New("not possible to plant seed")
	ErrInvalidPermissionsForSeed = errors.New("invalid permissions to retreive seed")
)

func (s *seedService) GetAllUserSeeds(ctx context.Context, userID string) ([]*SeedGroup, error) {
	userIDFromCtx, err := contextkeys.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if userID != userIDFromCtx {
		return nil, fmt.Errorf("you do not have authorised access")
	}

	seeds, err := s.store.Seed.GetAllByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	seedGroupsMap := make(map[string]*SeedGroup)
	for _, seed := range seeds {
		sg, ok := seedGroupsMap[seed.BotanicalName]
		if ok {
			sg.Seeds = append(sg.Seeds, seed)
			sg.Count++
			continue
		} else {
			seedGroupsMap[seed.BotanicalName] = &SeedGroup{
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

	targetCentre := models.Coordinates{Lat: dto.Latitude, Lng: dto.Longitude}
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

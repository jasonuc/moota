package services

import (
	"errors"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SeedService interface {
	GetAllUserSeeds(string) ([]*SeedGroup, error)
	GetSeed(string, string) (*models.Seed, error)
	PlantSeed(PlantSeedReqDto) (*models.Plant, error)
	WithStore(*store.Store) SeedService
}

type seedService struct {
	soilService  SoilService
	plantService PlantService
	store        *store.Store
}

func NewSeedService(store *store.Store) SeedService {
	return &seedService{
		store: store,
	}
}

func (s *seedService) WithStore(store *store.Store) SeedService {
	copy := *s
	copy.store = store
	return &copy
}

type PlantSeedReqDto struct {
	Longitude float64
	Latitude  float64
	SeedID    string
	UserID    string
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

func (s *seedService) GetAllUserSeeds(userID string) ([]*SeedGroup, error) {
	seeds, err := s.store.Seed.GetAllByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	seedGroups := make(map[string]*SeedGroup)
	for _, seed := range seeds {
		sg, ok := seedGroups[seed.BotanicalName]
		if ok {
			sg.Seeds = append(sg.Seeds, seed)
			sg.Count++
			continue
		} else {
			seedGroups[seed.BotanicalName] = &SeedGroup{
				BotanicalName: seed.BotanicalName,
				Count:         1,
				Seeds:         []*models.Seed{seed},
			}
		}
	}

	result := make([]*SeedGroup, 0, len(seedGroups))
	for _, group := range seedGroups {
		result = append(result, group)
	}

	return result, nil
}

func (s *seedService) GetSeed(userID, seedID string) (*models.Seed, error) {
	seed, err := s.store.Seed.Get(seedID)
	if err != nil {
		return nil, err
	}
	if seed.OwnerID == userID {
		return nil, ErrInvalidPermissionsForSeed
	}
	return seed, nil
}

func (s *seedService) PlantSeed(dto PlantSeedReqDto) (*models.Plant, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)
	soilServiceWithTx := s.soilService.WithStore(tx)
	plantServiceWithTx := s.plantService.WithStore(tx)

	seed, err := tx.Seed.Get(dto.SeedID)
	if err != nil {
		return nil, err
	}

	if seed.OwnerID != dto.UserID {
		return nil, ErrUnauthorizedSeedPlanting
	}

	targetCentre := models.Coordinates{Lat: dto.Latitude, Lng: dto.Longitude}
	plantCircleMeta := models.NewCircleMeta(targetCentre, models.PlantInteractionRadius)

	nearbySoils, err := tx.Soil.GetAllInProximity(targetCentre, models.SoilRadiusMLarge)
	if err != nil {
		return nil, err
	}

	if len(nearbySoils) == 0 {
		soil, err := soilServiceWithTx.CreateSoil(targetCentre, nearbySoils)
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

	plant, err := plantServiceWithTx.CreatePlant(targetSoil, seed, targetCentre)
	if err != nil {
		return nil, err
	}

	if err := tx.Seed.MarkAsPlanted(seed.ID); err != nil {
		return nil, err
	}

	if err := transaction.Commit(); err != nil {
		return nil, err
	}

	return plant, nil
}

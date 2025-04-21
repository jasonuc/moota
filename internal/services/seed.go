package services

import (
	"errors"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SeedService struct {
	soilService  *SoilService
	plantService *PlantService
	store        *store.Store
}

func NewSeedService(store *store.Store) *SeedService {
	return &SeedService{
		store: store,
	}
}

type PlantSeedReqDto struct {
	Longitude float64
	Latitude  float64
	SeedID    string
	UserID    string
}

var (
	ErrUnauthorizedSeedPlanting  = errors.New("not authorised to plant this seed")
	ErrNotPossibleToPlantSeed    = errors.New("not possible to plant seed")
	ErrInvalidPermissionsForSeed = errors.New("invalid permissions to retreive seed")
)

func (s *SeedService) GetAllUserSeeds(userID string) ([]*models.Seed, error) {
	seeds, err := s.store.Seed.GetAllByOwnerID(userID)
	if err != nil {
		return nil, err
	}
	return seeds, nil
}

func (s *SeedService) GetSeed(userID, seedID string) (*models.Seed, error) {
	seed, err := s.store.Seed.Get(seedID)
	if err != nil {
		return nil, err
	}
	if seed.OwnerID == userID {
		return nil, ErrInvalidPermissionsForSeed
	}
	return seed, nil
}

func (s *SeedService) PlantSeed(dto PlantSeedReqDto) (*models.Plant, error) {
	transaction, err := s.store.Begin()
	if err != nil {
		return nil, store.ErrTransactionCouldNotStart
	}
	//nolint:errcheck
	defer transaction.Rollback()

	tx := s.store.WithTx(transaction)
	soilServiceWithTx := s.soilService.withStore(tx)
	plantServiceWithTx := s.plantService.withStore(tx)

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
	return plant, nil
}

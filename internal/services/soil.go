package services

import (
	"context"
	"errors"
	"math"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SoilService interface {
	CreateSoil(context.Context, models.Coordinates, []*models.Soil) (*models.Soil, error)
	WithStore(*store.Store) SoilService
}

type soilService struct {
	store *store.Store
}

func NewSoilSerivce(store *store.Store) SoilService {
	return &soilService{
		store: store,
	}
}

func (s *soilService) WithStore(store *store.Store) SoilService {
	copy := *s
	copy.store = store
	return &copy
}

var (
	ErrNoSoilGenerated = errors.New("no soil generated")
)

func (s *soilService) CreateSoil(ctx context.Context, centre models.Coordinates, nearbySoils []*models.Soil) (*models.Soil, error) {
	radius := models.RandomSoilRadius(models.RandomSoilRadiusParam{MaxRadius: math.Inf(1)})
	newSoilCircleMeta := models.NewCircleMeta(centre, radius)
	soilMeta := models.RandomSoilMeta()

	overlappingSoils := make(map[string]*models.Soil)
	for _, soil := range nearbySoils {
		if newSoilCircleMeta.OverlapsWith(soil) {
			overlappingSoils[soil.ID] = soil
		}
	}

	if len(overlappingSoils) == 0 {
		newSoil := models.MapToNewSizedSoilFn(radius)(soilMeta, centre)
		err := s.store.Soil.Insert(ctx, newSoil)
		if err != nil {
			return nil, err
		}
		return newSoil, nil
	}

	filterForRadius := models.RandomSoilRadiusParam{MaxRadius: math.Inf(1)}
	for _, soil := range overlappingSoils {
		filterForRadius.MaxRadius = math.Min(filterForRadius.MaxRadius, s.maxSoilRadius(newSoilCircleMeta, soil))
	}

	radius = models.RandomSoilRadius(filterForRadius)
	if radius == models.SoilRadiusMZero {
		return nil, ErrNoSoilGenerated
	}

	soil := models.MapToNewSizedSoilFn(radius)(soilMeta, centre)
	err := s.store.Soil.Insert(ctx, soil)
	if err != nil {
		return nil, err
	}

	return soil, nil
}

func (s *soilService) maxSoilRadius(circleMeta models.CircleMeta, nearbySoil *models.Soil) float64 {
	d := circleMeta.Centre().DistanceM(nearbySoil.Centre())
	maxRadius := d - nearbySoil.RadiusM() - 0.1 // the reason for the 0.1 subtraction is to prevent a possible case of tangential soils
	return math.Max(0.0, maxRadius)
}

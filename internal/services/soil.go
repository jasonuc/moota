package services

import (
	"errors"
	"math"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type SoilService struct {
	store *store.Store
}

func NewSoilSerivce(store *store.Store) *SoilService {
	return &SoilService{
		store: store,
	}
}

func (s *SoilService) withStore(store *store.Store) *SoilService {
	copy := *s
	copy.store = store
	return &copy
}

var (
	ErrNoSoilGenerated = errors.New("no soil generated")
)

func (s *SoilService) CreateSoil(centre models.Coordinates, nearbySoils []*models.Soil) (*models.Soil, error) {
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
		err := s.store.Soil.Insert(newSoil)
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
	err := s.store.Soil.Insert(soil)
	if err != nil {
		return nil, err
	}

	return soil, nil
}

func (s *SoilService) maxSoilRadius(circleMeta models.CircleMeta, nearbySoil *models.Soil) float64 {
	d := circleMeta.Centre().DistanceM(nearbySoil.Centre())
	maxRadius := d - nearbySoil.RadiusM() - 0.1 // the reason for the 0.1 subtraction is to prevent a possible case of tangential soils
	return math.Max(0.0, maxRadius)
}

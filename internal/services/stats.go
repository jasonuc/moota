package services

import (
	"context"

	"github.com/jasonuc/moota/internal/models"
	"github.com/jasonuc/moota/internal/store"
)

type StatsService interface {
	GetBasicStats(ctx context.Context) (*models.Stats, error)
}

type statsService struct {
	store *store.Store
}

func NewStatsService(store *store.Store) StatsService {
	return &statsService{store: store}
}

// GetBasicStats implements StatsService.
func (s *statsService) GetBasicStats(ctx context.Context) (*models.Stats, error) {
	var stat models.Stats
	plantCount, err := s.store.Plant.Count(ctx)
	if err != nil {
		return nil, err
	}
	stat.PlantCount = plantCount
	return &stat, nil
}

var _ StatsService = (*statsService)(nil)

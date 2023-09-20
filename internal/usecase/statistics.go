package usecase

import (
	"context"
	"fmt"
	"tm-backend-trainee-impl-clean-template/internal/entity"
)

// StatisticsUseCase -.
type StatisticsUseCase struct {
	repo StatisticsRepo
}

// New -.
func New(r StatisticsRepo) *StatisticsUseCase {
	return &StatisticsUseCase{
		repo: r,
	}
}

// Save -.
func (uc *StatisticsUseCase) Save(ctx context.Context, t entity.Metrics) error {
	err := uc.repo.InsertOrUpdate(ctx, t)
	if err != nil {
		return fmt.Errorf("StatisticsUseCase - Save - s.repo.InsertOrUpdate: %w", err)
	}

	return nil
}

// Clean -.
func (uc *StatisticsUseCase) Clear(ctx context.Context) error {
	err := uc.repo.Truncate(ctx)
	if err != nil {
		return fmt.Errorf("StatisticsUseCase - Clear - s.repo.Truncate: %w", err)
	}

	return nil
}

// Get -.
func (uc *StatisticsUseCase) Get(ctx context.Context, g entity.DoGetRequest) ([]entity.Statistics, error) {
	r, err := uc.repo.GetStatistics(ctx, g)
	if err != nil {
		return []entity.Statistics{}, fmt.Errorf("StatisticsUseCase - Get - s.repo.GetStatistics: %w", err)
	}
	return r, nil
}

// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"tm-backend-trainee-impl-clean-template/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Statistics -.
	Statistics interface {
		Save(context.Context, entity.Metrics) error
		Get(context.Context, entity.DoGetRequest) ([]entity.Statistics, error)
		Clear(context.Context) error
	}

	// StatisticsRepo -.
	StatisticsRepo interface {
		InsertOrUpdate(context.Context, entity.Metrics) error
		GetStatistics(context.Context, entity.DoGetRequest) ([]entity.Statistics, error)
		Truncate(context.Context) error
	}
)

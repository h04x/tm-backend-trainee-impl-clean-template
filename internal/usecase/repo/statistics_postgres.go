package repo

import (
	"context"
	"fmt"

	"tm-backend-trainee-impl-clean-template/internal/entity"
	"tm-backend-trainee-impl-clean-template/pkg/postgres"
)

const _defaultEntityCap = 64

// PostgresRepo -.
type StatisticsRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *StatisticsRepo {
	return &StatisticsRepo{pg}
}

// InsertOrUpdate -.
func (r *StatisticsRepo) InsertOrUpdate(ctx context.Context, t entity.Metrics) error {
	sql, args, err := r.Builder.
		Insert("clicks").
		Columns("date, views, clicks, cost").
		Values(t.Date, t.Views, t.Clicks, t.Cost).
		Suffix("on conflict (date) do update set views = clicks.views + ?, clicks = clicks.clicks + ?, cost = clicks.cost + ?",
			t.Views, t.Clicks, t.Cost).
		ToSql()
	if err != nil {
		return fmt.Errorf("StatisticsRepo - InsertOrUpdate - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("StatisticsRepo - InsertOrUpdate - r.Pool.Exec: %w", err)
	}

	return nil
}

// Truncate -.
func (r *StatisticsRepo) Truncate(ctx context.Context) error {
	_, err := r.Pool.Exec(ctx, "truncate clicks")
	if err != nil {
		return fmt.Errorf("StatisticsRepo - Truncate - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetStatistics -.
func (r *StatisticsRepo) GetStatistics(ctx context.Context, t entity.DoGetRequest) ([]entity.Statistics, error) {
	sql, args, err := r.Builder.
		Select(`date::text, views, clicks, cost,
			round(coalesce(cost / NULLIF(clicks, 0), 0), 2) as cpc, 
			round(coalesce(cost / NULLIF(views, 0) * 1000, 0), 2) as cpm`).
		From("clicks").
		Where("date between ? and ?", t.From, t.To).
		OrderBy(t.Order).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StatisticsRepo - InsertOrUpdate - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("StatisticsRepo - GetStatistics - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Statistics, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.Statistics{}

		err = rows.Scan(&e.Metrics.Date, &e.Metrics.Views,
			&e.Metrics.Clicks, &e.Metrics.Cost, &e.Cpc, &e.Cpm)
		if err != nil {
			return nil, fmt.Errorf("StatisticsRepo - GetStatistics - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil

}

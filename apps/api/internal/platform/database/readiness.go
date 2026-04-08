package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadinessProbe struct {
	pool *pgxpool.Pool
}

func NewReadinessProbe(pool *pgxpool.Pool) *ReadinessProbe {
	return &ReadinessProbe{pool: pool}
}

func (probe *ReadinessProbe) Name() string {
	return "postgres"
}

func (probe *ReadinessProbe) Check(ctx context.Context) error {
	return probe.pool.Ping(ctx)
}

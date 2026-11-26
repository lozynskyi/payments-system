// Package postgres connect service to PostgreSQL through the pgx connection pool.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a connection pool for the database at dsn, and check that the database is reachable before
// returning it.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres parse dsn: %w", err)
	}

	// Pool sizing: a few warm connections, and a hard cap well below
	// Postgres's max_connections, so many replicas cans share one server.
	cfg.MinConns = 2
	cfg.MaxConns = 10

	// Recycle connections now and then, so the load rebalances after failover.
	// The jitter stops every connection from reconnecting at the same time.
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnLifetimeJitter = 5 * time.Minute
	cfg.MaxConnIdleTime = 15 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: create poll %w", err)
	}

	// The pool connects lazily. Ping now so the bad dsn or database that is down
	// fails at startup, instead of on the first request.

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("postgres: ping pool: %w", err)
	}
	return pool, nil
}

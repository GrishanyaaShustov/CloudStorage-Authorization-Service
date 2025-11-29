package postgres

import (
	"authorization-service/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New creates a new Postgres connection pool based on application config.
// It performs a ping to verify connectivity before returning the pool.
func New(ctx context.Context, log *slog.Logger, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	log.Info("connecting to Postgres",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("db", cfg.Name),
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Postgres DSN: %w", err)
	}

	poolCfg.MaxConns = 24
	poolCfg.MinConns = 1
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	// Open pool
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres pool: %w", err)
	}
	// Check, connection to db.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}
	log.Info("connected to Postgres successfully")

	return pool, nil
}

// MustNew creates a Postgres pool and panics on any error.
// Used in the application's startup layer (main.go).
func MustNew(ctx context.Context, log *slog.Logger, cfg config.DatabaseConfig) *pgxpool.Pool {
	pool, err := New(ctx, log, cfg)
	if err != nil {
		panic(err)
	}
	return pool
}

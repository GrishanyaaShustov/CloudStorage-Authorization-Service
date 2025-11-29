package redis

import (
	"authorization-service/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// New creates a new Redis client and verifies the connection via Ping.
// It does NOT panic — errors are returned to the caller.
func New(ctx context.Context, log *slog.Logger, cfg config.RedisConfig) (*goredis.Client, error) {

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	log.Info("connecting to Redis",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
	)

	rdb := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     20,
	})

	// Check connectivity
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info("connected to Redis successfully")
	return rdb, nil
}

// MustNew creates a Redis client and PANICS on error.
// This is used at application startup (in main.go).
func MustNew(ctx context.Context, log *slog.Logger, cfg config.RedisConfig) *goredis.Client {
	rdb, err := New(ctx, log, cfg)
	if err != nil {
		panic(err)
	}
	return rdb
}

package app

import (
	pgstorage "authorization-service/internal/storage/postgres"
	"context"
	"log/slog"

	grpcapp "authorization-service/internal/app/grpc"
	"authorization-service/internal/config"
)

// App is a top-level application container.
// It wires configuration, logger and sub-apps
type App struct {
	GRPC *grpcapp.App
}

// New builds the whole application graph.
// For now it only wires gRPC server
func New(log *slog.Logger, cfg *config.Config) *App {
	if log == nil {
		log = slog.Default()
	}

	// Adjust this according to your Config structure.
	// Here we assume cfg.GRPC has Port field of type int.
	grpcPort := cfg.GRPC.Port

	ctx := context.Background()
	pg := pgstorage.MustNew(ctx, log, cfg.Database)

	grpcApp := grpcapp.New(log, grpcPort, pg)

	return &App{
		GRPC: grpcApp,
	}
}

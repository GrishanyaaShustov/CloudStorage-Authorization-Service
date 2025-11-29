package main

import (
	"authorization-service/internal/app"
	"authorization-service/internal/config"
)

func main() {
	// 1. Init cfg
	cfg := config.MustLoad()

	// 2. Init logger
	logger := setupLogger(cfg.Env)

	// 3. Build application: wire config, logger, gRPC app, services, etc.
	application := app.New(logger, cfg)

	// 4. Run gRPC server (panic if cant start).
	application.GRPC.MustRun()

}

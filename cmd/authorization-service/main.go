package main

import (
	"authorization-service/internal/config"
	"authorization-service/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

func main() {
	cfg, _ := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.Any("config", cfg))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "prod":
		log = setupPrettySlog()
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

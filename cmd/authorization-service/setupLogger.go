package main

import (
	"authorization-service/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

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

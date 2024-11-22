package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var env string

func SetEnv(environment string) {
	env = environment
}

func Logger() *slog.Logger {
	switch env {
	case "local", "test":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "stage":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic(fmt.Sprintf("Unknown environment: %s", env))
	}
}

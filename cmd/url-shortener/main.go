package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"vigilant-octo-spoon/internal/config"
	middlewarelogger "vigilant-octo-spoon/internal/http_server/middlewares/logger"
	"vigilant-octo-spoon/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoadEnv()
	fmt.Println(cfg)
	logger := NewLogger(cfg.Env)
	logger.Info("starting server", slog.String("env config", cfg.Env))
	logger.Debug("Debug messages")

	storage, err := sqlite.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error("error creating storage", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middlewarelogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

}

const (
	envLocal = "local"
	envDev   = "development"
	evnProd  = "production"
)

func NewLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case evnProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}

package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"vigilant-octo-spoon/internal/config"
	"vigilant-octo-spoon/internal/http_server/handlers/alias/save"
	middlewarelogger "vigilant-octo-spoon/internal/http_server/middlewares/logger"
	"vigilant-octo-spoon/internal/storage/sqlite"
	"vigilant-octo-spoon/lib/logger/sl"
)

func main() {
	cfg := config.MustLoadEnv()
	fmt.Println(cfg)
	logger := NewLogger(cfg.Env)
	logger.Info("starting server", slog.String("env config", cfg.Env))
	logger.Debug("Debug messages")

	storage, err := sqlite.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error("error creating storage", sl.Err(err))
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middlewarelogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, storage))

	logger.Info("starting server on address", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Error("error starting server", sl.Err(err))
	}
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

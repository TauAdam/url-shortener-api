package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	grpcClient "github.com/tauadam/url-shortener-api/internal/clients/sso/grpc"
	"github.com/tauadam/url-shortener-api/internal/config"
	deleteURL "github.com/tauadam/url-shortener-api/internal/http_server/handlers/alias/delete"
	"github.com/tauadam/url-shortener-api/internal/http_server/handlers/alias/save"
	"github.com/tauadam/url-shortener-api/internal/http_server/handlers/redirect"
	middlewarelogger "github.com/tauadam/url-shortener-api/internal/http_server/middlewares/logger"
	"github.com/tauadam/url-shortener-api/internal/storage/sqlite"
	"github.com/tauadam/url-shortener-api/lib/logger/handler/pretty_slog"
	"github.com/tauadam/url-shortener-api/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoadEnv()
	fmt.Println(cfg)

	logger := NewLogger(cfg.Env)

	logger.Info("starting server", slog.String("env config", cfg.Env))
	logger.Debug("Debug messages")

	_, err := grpcClient.New(
		context.Background(),
		logger,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.RetriesNumber,
		cfg.Clients.SSO.Timeout,
	)
	if err != nil {
		logger.Error("error creating sso client", sl.Err(err))

	}

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

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth(
			"url-shortener-api",
			map[string]string{cfg.HttpServerConfig.User: cfg.HttpServerConfig.Password},
		))
		r.Post("/", save.New(logger, storage))
		r.Delete("/{alias}", deleteURL.New(logger, storage))
	})
	router.Get("/{alias}", redirect.New(logger, storage))

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
		logger = setupPrettyLogger()
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case evnProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}

func setupPrettyLogger() *slog.Logger {
	options := pretty_slog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := options.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

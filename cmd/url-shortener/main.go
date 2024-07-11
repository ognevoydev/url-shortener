package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	mwLogger "url-shortener/internal/api/middleware/logger"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	cfg := config.Load()

	log := logger.Setup(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatal("failed to init storage", slog.String("error", err.Error()))
	}
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))
}

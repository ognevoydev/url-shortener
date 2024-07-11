package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	save "url-shortener/internal/api/handlers/url"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}

package main

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/api/routes"
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

	router := routes.Setup(log, storage)

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

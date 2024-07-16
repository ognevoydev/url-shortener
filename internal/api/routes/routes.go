package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"url-shortener/internal/api/handlers/delete"
	"url-shortener/internal/api/handlers/redirect"
	"url-shortener/internal/api/handlers/register"
	"url-shortener/internal/api/handlers/save"
	mwLogger "url-shortener/internal/api/middleware/logger"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/storage/sqlite"
)

func Setup(log *logger.Logger, storage *sqlite.Storage) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))

	// URLs
	router.Post("/save", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/{alias}", delete.New(log, storage))

	// Users
	router.Post("/register", register.New(log, storage))

	return router
}

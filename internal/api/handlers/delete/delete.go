package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/api/response"
	"url-shortener/internal/lib/logger"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLRemover
type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *logger.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.delete.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		err := urlRemover.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("delete url error"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		// TODO - return error if url not exists

		render.JSON(w, r, response.OK())
	}
}

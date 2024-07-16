package register

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/api/response"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/storage"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserCreator interface {
	CreateUser(username string, password string) (int64, error)
}

func New(log *logger.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.register.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		id, err := userCreator.CreateUser(req.Username, req.Password)

		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Info("username already exists", slog.String("username", req.Username))
				render.JSON(w, r, response.Error("username already exists"))
				return
			} else {
				log.Info("failed to create user", slog.String("username", req.Username))
				render.JSON(w, r, response.Error("failed to create user"))
				return
			}
		}

		log.Info("user created", slog.Int64("id", id))

		render.JSON(w, r, response.OK())
	}
}

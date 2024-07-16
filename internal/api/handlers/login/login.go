package login

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/api/response"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/lib/security"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
	Token string `json:"token"`
}

type UserAuthenticator interface {
	AuthenticateUser(username string, password string) (int64, error)
	CreateSession(userId int64, token string) (int64, error)
}

func New(log *logger.Logger, authenticator UserAuthenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.login.New"

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

		userId, err := authenticator.AuthenticateUser(req.Username, req.Password)
		if userId == 0 || err != nil {
			log.Info("failed to login user", slog.String("username", req.Username))
			render.JSON(w, r, response.Error("authentication failed"))
			return
		}

		log.Info("user authenticated", slog.String("username", req.Username))

		token := security.GenerateToken()

		_, err = authenticator.CreateSession(userId, token)
		if err != nil {
			log.Info("failed to create session", slog.String("username", req.Username))
			render.JSON(w, r, response.Error("failed to create session"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Token:    token,
		})
	}
}

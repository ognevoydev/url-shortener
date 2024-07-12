package handlers

import (
	"context"
	"log/slog"
	"url-shortener/internal/lib/logger"
)

func NewDiscardLogger() *logger.Logger {
	return &logger.Logger{Logger: slog.New(&DiscardHandler{})}
}

type DiscardHandler struct{}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

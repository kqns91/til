package main

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const (
	requestIDKey = contextKey("requestID")
)

type Handler struct {
	slog.Handler
}

func NewHandler(h slog.Handler) *Handler {
	return &Handler{Handler: h}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String(string(requestIDKey), ctx.Value(requestIDKey).(string)))
	return h.Handler.Handle(ctx, r)
}

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, requestIDKey, "12345")

	logger := slog.New(NewHandler(slog.NewJSONHandler(os.Stdout, nil)))
	slog.SetDefault(logger)

	logger.InfoContext(ctx, "hoge")
}

package app_context

import (
	"context"
	"log/slog"
)

type contextKey int

const (
	_ contextKey = iota
	userIDKey
	loggerKey
)

func ContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)

	return userID, ok
}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		logger = slog.Default()
	}
	return logger
}

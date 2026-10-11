package middleware

import (
	"context"
	"log/slog"
	"transcoder/internal/presentation/messaging"
)

type fallbackFunc func(ctx context.Context, msg messaging.Message) error

func OnFailure(callback fallbackFunc) messaging.Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return messaging.HandlerFunc(func(ctx context.Context, msg messaging.Message) error {
			err := next.Handle(ctx, msg)
			if err == nil {
				return nil
			}

			slog.Error("error while handler message", "error", err)

			return callback(ctx, msg)
		})
	}
}

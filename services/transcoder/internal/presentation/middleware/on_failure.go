package middleware

import (
	"context"
	"transcoder/internal/presentation/messaging"
)

type fallbackFunc func(ctx context.Context, msg messaging.Message) error

func OnFailure(fallback fallbackFunc) messaging.Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return messaging.HandlerFunc(func(ctx context.Context, msg messaging.Message) error {
			err := next.Handle(ctx, msg)
			if err == nil {
				return nil
			}

			return fallback(ctx, msg)
		})
	}
}

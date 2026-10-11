package middleware

import (
	"context"
	"transcoder/internal/presentation/messaging"
)

type successFunc func(ctx context.Context, msg messaging.Message) error

func OnSuccess(callback successFunc) messaging.Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return messaging.HandlerFunc(func(ctx context.Context, msg messaging.Message) error {
			err := next.Handle(ctx, msg)
			if err != nil {
				return err
			}

			return callback(ctx, msg)
		})
	}
}

package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"transcoder/internal/presentation/messaging"
)

func Retry(attempts int, delay time.Duration) messaging.Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return messaging.HandlerFunc(func(ctx context.Context, msg messaging.Message) error {
			for attempt := 0; attempt < attempts; attempt += 1 {
				err := next.Handle(ctx, msg)
				if err == nil {
					return nil
				}

				slog.Warn("handler failed, retrying", "type", msg.Type, "attempt", attempt+1)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}
			}

			return fmt.Errorf("error after retrying handler")
		})
	}
}

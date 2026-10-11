package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	app_ports "transcoder/internal/application/ports"
	"transcoder/internal/presentation/messaging"
)

func Retry(attempts int, delay time.Duration) messaging.Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return messaging.HandlerFunc(func(ctx context.Context, msg messaging.Message) error {
			var err error
			for attempt := 0; attempt < attempts; attempt += 1 {
				err = next.Handle(ctx, msg)
				if err == nil || errors.Is(err, app_ports.ErrUnprocessable) {
					return err
				}

				slog.Warn("handler failed, retrying", "type", msg.Type, "attempt", attempt+1)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}
			}

			return fmt.Errorf("error after retrying handler: %w", err)
		})
	}
}

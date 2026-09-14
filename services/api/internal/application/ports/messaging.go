package app_ports

import (
	"context"
	"project/internal/domain/event"
)

type Publisher interface {
	PublicEvents(ctx context.Context, events []event.Interface) error
}

package pres_ports

import (
	"context"
	domain_event "transcoder/internal/domain/event"
)

type Publisher interface {
	PublishEvents(ctx context.Context, events []domain_event.Interface) error
}

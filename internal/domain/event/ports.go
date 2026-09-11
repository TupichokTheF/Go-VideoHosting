package event

import "time"

type Interface interface {
	EventName() string
	Payload() map[string]any
	OccurredAt() time.Time
}

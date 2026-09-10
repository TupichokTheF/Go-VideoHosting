package event

import "time"

type Interface interface {
	EventName() string
	Payload() map[string]any
	OccuredAt() time.Time
}

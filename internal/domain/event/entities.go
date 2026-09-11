package event

import "time"

type Base struct {
	CreatedAt time.Time
}

func (b Base) OccurredAt() time.Time {
	return b.CreatedAt
}

type Recorder struct {
	events []Interface
}

func (r *Recorder) Add(event Interface) {
	r.events = append(r.events, event)
}

func (r *Recorder) Pull() []Interface {
	events := r.events
	r.events = nil

	return events
}

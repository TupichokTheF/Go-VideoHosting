package event


type Base struct {
	CreatedAt int64
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
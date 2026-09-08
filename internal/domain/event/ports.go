package event


type Interface interface {
	EventName() string
	Payload() map[string]any
}
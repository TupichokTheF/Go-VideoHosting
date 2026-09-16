package app_ports

import "context"

type Message struct {
	Type     string
	Payload  map[string]any
	Callback CommitMessage
}

type CommitMessage func(ctx context.Context) error

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}

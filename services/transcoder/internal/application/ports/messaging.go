package app_ports

import "context"

type Message struct {
	Type    string
	Payload []byte
}

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}

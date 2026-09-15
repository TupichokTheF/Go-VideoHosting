package app_ports

import "context"


type Message struct {
	Type string
	Payload []byte
}

type Handler func(ctx context.Context, msg Message) error
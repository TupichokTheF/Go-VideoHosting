package messaging


type Bus struct {
	consumer consumerInterface
	handlers map[string]handlerInterface
}

type consumerInterface interface {
	FetchMessage() error
}

type handlerInterface interface {
	Handle() error
}

func NewBus(consumer consumerInterface, handlers map[string]handlerInterface) *Bus {
	return &Bus{
		consumer: consumer,
		handlers: handlers,
	}
}

func (bus *Bus) ProcessMessage() error {
	return nil
}
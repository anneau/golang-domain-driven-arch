package port

import "context"

type Message interface {
	Body() []byte
	Ack(ctx context.Context) error
	Nack(ctx context.Context) error
	MessageID() string
}

type EventSubscriber interface {
	Subscribe(ctx context.Context) (<-chan Message, error)
	Close() error
}

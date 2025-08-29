package messaging

import (
	"context"
)

// Message — сообщение
type Message struct {
	Key   []byte
	Value []byte
}

type Producer interface {
	ProduceMessage(ctx context.Context, topic string, msg Message) error
	Close() error
}

type Handler func(ctx context.Context, msg Message) error

type Consumer interface {
	ConsumeMessage(ctx context.Context, handler Handler) error
	Close() error
}

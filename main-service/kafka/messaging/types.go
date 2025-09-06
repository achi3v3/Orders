package messaging

import (
	"context"

	"github.com/segmentio/kafka-go"
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
	Run(ctx context.Context)
	ConsumeMessage(ctx context.Context) error
	Commit(ctx context.Context, msg kafka.Message) error
	Close() error
}

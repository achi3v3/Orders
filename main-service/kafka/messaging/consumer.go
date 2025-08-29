package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{reader: reader}
}

func (c *KafkaConsumer) ConsumeMessage(ctx context.Context, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			kafkaMsg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("KafkaConsumer.Consume: failed to read message: %w", err)
			}
			msg := Message{
				Key:   kafkaMsg.Key,
				Value: kafkaMsg.Value,
			}
			if err := handler(ctx, msg); err != nil {
				log.Printf("KafkaConsumer.Consume: failed to read message: %v", err)
			}
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

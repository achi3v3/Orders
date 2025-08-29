package messaging

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
		BatchTimeout:           1 * time.Second,
		WriteTimeout:           10 * time.Second,
	}
	return &kafkaProducer{writer: writer}
}

func (p *kafkaProducer) ProduceMessage(ctx context.Context, topic string, msg Message) error {
	kafkaMsg := kafka.Message{
		Topic: topic,
		Key:   msg.Key,
		Value: msg.Value,
	}
	return p.writer.WriteMessages(ctx, kafkaMsg)
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}

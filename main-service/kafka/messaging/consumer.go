package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"orders/internal/models"
	"orders/internal/subs"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	logger  *logrus.Logger
	handler *subs.Handler
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, logger *logrus.Logger, handler *subs.Handler) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		reader:  reader,
		logger:  logger,
		handler: handler,
	}
}

func (c *KafkaConsumer) Run(ctx context.Context) {
	c.logger.Info("KafkaConsumer: Starting consumer...")
	c.logger.Infof("KafkaConsumer: Brokers: %v, Topic: %s, GroupID: %s",
		c.reader.Config().Brokers,
		c.reader.Config().Topic,
		c.reader.Config().GroupID)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("KafkaConsume.Run: Consumer stop (ctx cancel)")
			return
		default:
			if err := c.ConsumeMessage(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				c.logger.Errorf("KafkaConsumer: Error consuming message: %v", err)
				time.Sleep(1 * time.Second)
			}
		}
	}

}
func (c *KafkaConsumer) ConsumeMessage(ctx context.Context) error {
	kafkaMsg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		c.logger.Errorf("KafkaConsumer.Consume: %v", err)
		return fmt.Errorf("KafkaConsumer.Consume: %w", err)
	}
	c.logger.Infof("KafkaConsumer.runConsumer: Получено: Key: %s | Value: %s", string(kafkaMsg.Key), string(kafkaMsg.Value))
	var order *models.OrderJson
	if err := json.Unmarshal(kafkaMsg.Value, &order); err != nil {
		c.logger.Errorf("consumer.runConsumer: %v", err)
		return fmt.Errorf("KafkaConsumer: %v", err)
	}
	log := c.logger.WithField("order_uid", order.OrderUID)
	log.Info("Read Message")
	if err = c.handler.Create(ctx, order); err != nil {
		log.Errorf("KafkaConsumer.ConsumeMessage: %v", err)
		return fmt.Errorf("KafkaConsumer.ConsumeMessage: %w", err)

	}
	if err = c.Commit(ctx, kafkaMsg); err != nil {
		log.Errorf("KafkaConsumer.ConsumeMessage: %v", err)
		return fmt.Errorf("KafkaConsumer.ConsumeMessage: %w", err)

	}
	log.Info("Commit Message")
	return nil
}

func (c *KafkaConsumer) Commit(ctx context.Context, msg kafka.Message) error {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

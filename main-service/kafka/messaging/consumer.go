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
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})

	return &KafkaConsumer{
		reader:  reader,
		logger:  logger,
		handler: handler,
	}
}

func (c *KafkaConsumer) Run(ctx context.Context) {
	c.logger.Info("KafkaConsumer.Run: Starting consumer...")
	c.logger.Infof("KafkaConsumer.Run: Brokers: %v, Topic: %s, GroupID: %s",
		c.reader.Config().Brokers,
		c.reader.Config().Topic,
		c.reader.Config().GroupID)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("KafkaConsume.Run: Consumer stop (context canceled)")
			return
		default:
			if err := c.ConsumeMessage(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				c.logger.Errorf("KafkaConsumer: Error consuming message: %v", err)
				time.Sleep(2 * time.Second)
			}
		}
	}

}
func (c *KafkaConsumer) ConsumeMessage(ctx context.Context) error {
	kafkaMsg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		c.logger.Errorf("KafkaConsumer.ConsumeMessage: failed to fetch msg: %v", err)
		return fmt.Errorf("fetch message: %w", err)
	}
	log := c.logger.WithFields(logrus.Fields{
		"topic":     kafkaMsg.Topic,
		"partition": kafkaMsg.Partition,
		"offset":    kafkaMsg.Offset,
		"key":       kafkaMsg.Key,
	})
	log.Info("KafkaConsumer.Run: Received kafka message")
	var order *models.OrderJson
	if err := json.Unmarshal(kafkaMsg.Value, &order); err != nil {
		log.Errorf("KafkaConsumer.ConsumeMessage: Failed to unmarshal message: %v. Message %s", err, string(kafkaMsg.Value))

		if commitErr := c.reader.CommitMessages(ctx, kafkaMsg); commitErr != nil {
			log.Errorf("KafkaConsumer.ConsumeMessage: Failed to commit invalid message: %v", commitErr)
			return fmt.Errorf("commit invalid message: %w", commitErr)
		}
		log.Warn("KafkaConsumer.ConsumeMessage: Invalid message skipped and committed")
		return nil
	}

	log = log.WithField("order_uid", order.OrderUID)
	log.Info("KafkaConsumer.Run: Read Message")
	if err = c.handler.Create(ctx, order); err != nil {
		log.Errorf("KafkaConsumer.ConsumeMessage: %v", err)
		return fmt.Errorf("failed to create order %s. error: %w", order.OrderUID, err)

	}
	if err = c.Commit(ctx, kafkaMsg); err != nil {
		log.Errorf("KafkaConsumer.ConsumeMessage: %v", err)
		return fmt.Errorf("failed to commit message %s. error: %w", order.OrderUID, err)

	}
	log.Info("KafkaConsumer.Run: Commit Message")
	return nil
}

func (c *KafkaConsumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

func (c *KafkaConsumer) Close() error {
	c.logger.Info("KafkaConsumer.Close: Closing Kafka consumer")
	return c.reader.Close()
}

package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"orders/internal/models"
	"orders/internal/subs"
	"orders/kafka/messaging"
)

func RunConsumer(ctx context.Context, consumer messaging.Consumer, handl *subs.Handler) {
	handler := func(ctx context.Context, msg messaging.Message) error {
		log.Printf("kafka.runConsumer: Получено: Key: %s | Value: %s", string(msg.Key), string(msg.Value))
		data := &models.OrderJson{}
		if err := json.Unmarshal(msg.Value, data); err != nil {
			return fmt.Errorf("kafka.RunConsumer: %w", err)
		}
		handl.Create(data)
		return nil
	}
	if err := consumer.ConsumeMessage(ctx, handler); err != nil {
		log.Fatalf("kafka.runConsumer: Ошибка консьюмера")
	}

}

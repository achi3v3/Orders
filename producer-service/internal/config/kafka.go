package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	KafkaURL      string
	Topic         string
	GroupConsumer string
}

func LoadKafkaConfig() (*KafkaConfig, error) {
	envPath := filepath.Join("..", "configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("config.LoadPostgresConfig: %w", err)
	}
	config := &KafkaConfig{
		KafkaURL:      getEnv("KAFKA_URL", "kafka:9092"),
		Topic:         getEnv("TEST_TOPIC", "test_topic"),
		GroupConsumer: getEnv("GROUP_ID", "test_group"),
	}
	return config, nil
}

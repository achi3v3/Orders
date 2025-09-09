package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type KafkaConfig struct {
	KafkaURL      string
	Topic         string
	GroupConsumer string
	Logger        *logrus.Logger
}

func LoadKafkaConfig(logger *logrus.Logger) (*KafkaConfig, error) {
	envPath := filepath.Join("configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		logger.Errorf("config.LoadPostgresConfig: %v", err)
		return nil, fmt.Errorf("error with load env: %w", err)
	}
	config := &KafkaConfig{
		KafkaURL:      GetEnv("KAFKA_URL", "kafka:9092"),
		Topic:         GetEnv("TEST_TOPIC", "test_topic"),
		GroupConsumer: GetEnv("GROUP_ID", "test_group"),
		Logger:        logger,
	}
	return config, nil
}

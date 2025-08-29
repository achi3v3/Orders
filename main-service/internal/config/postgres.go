package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadPostgresConfig() (*PostgresConfig, error) {
	envPath := filepath.Join("..", "configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("config.LoadPostgresConfig: %w", err)
	}
	config := &PostgresConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_NAME", "Orders"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
	}
	return config, nil
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

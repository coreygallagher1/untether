package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Kafka
	KafkaBrokers []string
	KafkaGroupID string

	// Server
	GRPCPort int
	HTTPPort int

	// Plaid
	PlaidClientID     string
	PlaidSecret      string
	PlaidEnvironment string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		// Database defaults
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvAsIntOrDefault("DB_PORT", 5432),
		DBUser:     getEnvOrDefault("DB_USER", "untether"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "untether"),
		DBName:     getEnvOrDefault("DB_NAME", "untether"),

		// Kafka defaults
		KafkaBrokers: []string{getEnvOrDefault("KAFKA_BROKERS", "localhost:9092")},
		KafkaGroupID: getEnvOrDefault("KAFKA_GROUP_ID", "untether-group"),

		// Server defaults
		GRPCPort: getEnvAsIntOrDefault("GRPC_PORT", 50051),
		HTTPPort: getEnvAsIntOrDefault("HTTP_PORT", 8080),

		// Plaid (required)
		PlaidClientID:     getEnvOrDefault("PLAID_CLIENT_ID", ""),
		PlaidSecret:      getEnvOrDefault("PLAID_SECRET", ""),
		PlaidEnvironment: getEnvOrDefault("PLAID_ENVIRONMENT", "sandbox"),
	}

	// Validate required fields
	if config.PlaidClientID == "" || config.PlaidSecret == "" {
		return nil, fmt.Errorf("PLAID_CLIENT_ID and PLAID_SECRET are required")
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
} 
package config

import "os"

type Config struct {
	Port                       string
	DatabaseURL                string
	RabbitMQURL                string
	RabbitMQExchange           string
	RabbitMQQueue              string
	RabbitMQRoutingKey         string
	RabbitMQDeadLetterExchange string
	RabbitMQDeadLetterQueue    string
}

func Load() Config {
	return Config{
		Port:                       getEnv("PORT", "8090"),
		DatabaseURL:                os.Getenv("DATABASE_URL"),
		RabbitMQURL:                os.Getenv("RABBITMQ_URL"),
		RabbitMQExchange:           getEnv("RABBITMQ_EXCHANGE", "webhook.events"),
		RabbitMQQueue:              getEnv("RABBITMQ_QUEUE", "webhook.delivery"),
		RabbitMQRoutingKey:         getEnv("RABBITMQ_ROUTING_KEY", "webhook.created"),
		RabbitMQDeadLetterExchange: getEnv("RABBITMQ_DLX", "webhook.events.dlx"),
		RabbitMQDeadLetterQueue:    getEnv("RABBITMQ_DLQ", "webhook.delivery.dlq"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

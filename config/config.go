package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DbUrl    string
	Env      string
	LogLevel string
	// DigitalOcean Spaces Configuration
	DORegion      string
	DOEndpoint    string
	DOAccessKey   string
	DOSecretKey   string
	DOBucket      string
	PresignExpiry int // in minutes

	// redis
	RedisURL string

	RabbitMQ RabbitMQConfig
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Port:          getkey("PORT", "8080"),
		DbUrl:         getkey("DB_URL", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		Env:           getkey("ENV", "development"),
		LogLevel:      getkey("LOG_LEVEL", "info"),
		DORegion:      getkey("DO_REGION", "blr1"),
		DOEndpoint:    getkey("DO_ENDPOINT", "https://blr1.digitaloceanspaces.com"),
		DOAccessKey:   getkey("DO_ACCESS_KEY", ""),
		DOSecretKey:   getkey("DO_SECRET_KEY", ""),
		DOBucket:      getkey("DO_BUCKET", "my-bucket"),
		PresignExpiry: getEnvValue(getkey("PRESIGN_EXPIRY", "15"), 15),

		// redis
		RedisURL: getkey("REDIS_URL", ""),

		RabbitMQ: RabbitMQConfig{
			BrokerLink:   getkey("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			ExchangeName: getkey("EXCHANGE_NAME", "scarter.embedding"),
			ExchangeType: getkey("EXCHANGE_TYPE", "direct"),
			QueueName:    getkey("QUEUE_NAME", "source.processor.queue"),
			RoutingKey:   getkey("ROUTING_KEY", "source.process"),
			WorkerCount:  getEnvValue("WORKER_COUNT", 5),
		},
	}

}

func getkey(key string, fallback string) string {
	if os.Getenv(key) == "" {
		return fallback
	}
	return os.Getenv(key)
}

func getEnvValue(s string, fallback int) int {
	if s == "" {
		return fallback
	}

	value, err := strconv.Atoi(s)

	if err != nil {
		return fallback
	}
	return value
}

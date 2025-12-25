package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DbUrl       string
	Env         string
	JwtSecret   string
	LogLevel    string
	RabbitMQUrl string
	QueueName   string
	// S3 Configuration
	AWSRegion     string
	S3BucketName  string
	PresignExpiry int // in minutes

	// redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:          getkey("PORT", "8080"),
		DbUrl:         getkey("DB_URL", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		Env:           getkey("ENV", "development"),
		JwtSecret:     getkey("JWT_SECRET", "secret"),
		LogLevel:      getkey("LOG_LEVEL", "info"),
		RabbitMQUrl:   getkey("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"),
		QueueName:     getkey("QUEUE_NAME", "carter_queue"),
		AWSRegion:     getkey("AWS_REGION", "us-east-1"),
		S3BucketName:  getkey("S3_BUCKET_NAME", "carter-sources"),
		PresignExpiry: getEnvValue(getkey("PRESIGN_EXPIRY", "15"), 15),

		// redis
		RedisAddr:     getkey("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getkey("REDIS_PASSWORD", ""),
		RedisDB:       getEnvValue(getkey("REDIS_DB", "0"), 0),
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

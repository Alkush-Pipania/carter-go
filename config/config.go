package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DbUrl string
	Env   string
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:  getkey("PORT", "8080"),
		DbUrl: getkey("DB_URL", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		Env:   getkey("ENV", "development"),
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

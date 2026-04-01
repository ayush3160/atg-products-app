package app

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr      string
	MongoURI      string
	MongoDatabase string
}

func LoadConfig() Config {
	return Config{
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		MongoURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGODB_DATABASE", "sample_atg"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func ParseDurationEnv(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

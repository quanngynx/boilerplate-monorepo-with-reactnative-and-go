package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GinMode        string
	AllowedOrigins []string
	DatabaseURL    string
	RedisURL       string
	LogLevel       string
	AppName        string
}

func Load() (*Config, error) {
	loadEnvFiles()

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		AppName:     getEnv("APP_NAME", "server"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}

	origins := getEnv("ALLOWED_ORIGINS", "http://localhost:3001")
	cfg.AllowedOrigins = splitAndTrim(origins)

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}

	return cfg, nil
}

func loadEnvFiles() {
	_ = godotenv.Load(".env")

	if os.Getenv("DATABASE_URL") == "" || os.Getenv("REDIS_URL") == "" {
		_ = godotenv.Load(".env.example")
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

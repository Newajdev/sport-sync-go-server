package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Dsn           string
	JwtSecret     string
	AdminEmail    string
	AdminPassword string
	AdminName     string
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	return &Config{
		Port:          os.Getenv("PORT"),
		Dsn:           os.Getenv("DSN"),
		JwtSecret:     os.Getenv("JWT_SECRET"),
		AdminEmail:    envOrDefault("ADMIN_EMAIL", "admin@spotsync.com"),
		AdminPassword: envOrDefault("ADMIN_PASSWORD", "admin123456"),
		AdminName:     envOrDefault("ADMIN_NAME", "SpotSync Admin"),
	}
}

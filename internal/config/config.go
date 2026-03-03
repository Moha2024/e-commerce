package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN       string
	Port      string
	JWTSecret string
	RedisAddr string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	redisAddr := os.Getenv("REDIS_ADDR")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("JWT_SECRET") == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	return &Config{
		DSN:       dsn,
		Port:      port,
		JWTSecret: os.Getenv("JWT_SECRET"),
		RedisAddr: redisAddr,
	}, nil
}

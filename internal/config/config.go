package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_NAME string
	APP_ENV  string
	APP_PORT string

	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_URL      string
	DB_SSL_MODE string

	JWT_SECRET_KEY       string
	JWT_EXPIRATION time.Duration
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}

	return value
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return fallback
	}

	return value
}

func getRequireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}
func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed loading .env file: %w", err)
		}
	}

	jwtSecret, err := getRequireEnv("JWT_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	dbHost, err := getRequireEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbPort, err := getRequireEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbUser, err := getRequireEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPass, err := getRequireEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getRequireEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dbSSL, err := getRequireEnv("DB_SSL_MODE")
	if err != nil {
		return nil, err
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbSSL,
	)

	cfg := &Config{
		APP_NAME: getEnv("APP_NAME", "My App"),
		APP_ENV:  getEnv("APP_ENV", "development"),
		APP_PORT: getEnv("APP_PORT", "8080"),

		DB_HOST:     dbHost,
		DB_PORT:     dbPort,
		DB_USER:     dbUser,
		DB_PASSWORD: dbPass,
		DB_NAME:     dbName,
		DB_URL:      dbURL,
		DB_SSL_MODE: dbSSL,

		JWT_SECRET_KEY:       jwtSecret,
		JWT_EXPIRATION: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
	}

	return cfg, nil
}

package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv             string
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	JWTExpirationHours int
	MigrationsPath     string
}

func Load() (*Config, error) {
	jwtExpirationHours, err := getEnvAsInt("JWT_EXPIRATION_HOURS", 24)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", ""),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", ""),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", ""),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me"),
		JWTExpirationHours: jwtExpirationHours,
		MigrationsPath:     getEnv("MIGRATIONS_PATH", "file://migrations"),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required JWT_SECRET")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"task-api-go-ionix/internal/config"
)

func NewPostgresPool(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := BuildDatabaseURL(cfg)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func BuildDatabaseURL(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}

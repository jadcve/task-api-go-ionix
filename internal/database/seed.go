package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/security"
)

func SeedInitialAdmin(db *pgxpool.Pool) error {
	adminEmail := getEnv("ADMIN_EMAIL", "admin@test.com")
	adminPassword := getEnv("ADMIN_PASSWORD", "Admin123")
	adminName := getEnv("ADMIN_NAME", "Administrator")

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	if err := db.QueryRow(context.Background(), query, adminEmail).Scan(&exists); err != nil {
		return fmt.Errorf("check admin existence: %w", err)
	}

	if exists {
		return nil
	}

	hash, err := security.HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	insert := `
		INSERT INTO users (name, email, password_hash, role, must_change_password, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = db.Exec(
		context.Background(),
		insert,
		adminName,
		adminEmail,
		hash,
		string(domain.UserRoleAdmin),
		false,
		true,
	)
	if err != nil {
		return fmt.Errorf("insert admin user: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

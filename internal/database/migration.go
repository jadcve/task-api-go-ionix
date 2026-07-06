package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string, migrationsPath string) (err error) {
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if err == nil && sourceErr != nil {
			err = fmt.Errorf("close migration source: %w", sourceErr)
		}
		if err == nil && dbErr != nil {
			err = fmt.Errorf("close migration database: %w", dbErr)
		}
	}()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations up: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return nil
}

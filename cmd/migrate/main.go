package main

import (
	"log"
	"os"

	"task-api-go-ionix/internal/config"
	"task-api-go-ionix/internal/database"
)

func main() {
	action := "up"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	databaseURL := database.BuildDatabaseURL(cfg)

	switch action {
	case "up":
		if err := database.RunMigrations(databaseURL, cfg.MigrationsPath); err != nil {
			log.Fatalf("failed to run migrations up: %v", err)
		}
		log.Printf("migrations up executed successfully or no changes pending")
	case "down":
		log.Printf("migrate-down is not implemented from Go yet; use golang-migrate CLI in a later phase")
	default:
		log.Fatalf("invalid action %q. use: up | down", action)
	}
}

package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/4vertak/hestia/internal/config"
	"github.com/4vertak/hestia/internal/database"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate {up|down}")
		os.Exit(2)
	}
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	if os.Args[1] == "up" {
		err = database.RunMigrations(db)
	} else {
		err = database.RollbackMigrations(db)
	}
	if err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
}

package main

import (
	"log/slog"
	"os"

	"movie-platform/internal/config"
	"movie-platform/internal/database"
	"movie-platform/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration failed", "error", err)
		os.Exit(1)
	}
	db, err := database.Open(cfg)
	if err != nil {
		logger.Error("database failed", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("database pool failed", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
	if err := migrations.Up(db); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database migrations completed")
}

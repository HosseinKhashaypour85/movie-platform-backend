package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movie-platform/internal/config"
	"movie-platform/internal/database"
	"movie-platform/modules/media"
	"movie-platform/modules/transcoder"
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
	storage, err := media.NewS3PresignerFromEnv()
	if err != nil || storage == nil {
		logger.Error("S3 storage is required for transcoding", "error", err)
		os.Exit(1)
	}
	worker, err := transcoder.New(db, storage)
	if err != nil {
		logger.Error("transcoder initialization failed", "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger.Info("transcoder worker started")
	for {
		if ctx.Err() != nil {
			break
		}
		if err := worker.Run(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			logger.Error("transcode job failed", "error", err)
		}
		select {
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
		}
	}
	logger.Info("transcoder worker stopped")
}

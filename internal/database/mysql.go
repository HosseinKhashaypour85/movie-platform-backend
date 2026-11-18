package database

import (
	"context"
	"fmt"
	"time"

	"movie-platform/internal/config"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if cfg.Environment == "production" {
		dsn += "&tls=true"
	}
	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil { return nil, fmt.Errorf("open mysql: %w", err) }
	sqlDB, err := db.DB()
	if err != nil { return nil, fmt.Errorf("get mysql pool: %w", err) }
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil { _ = sqlDB.Close(); return nil, fmt.Errorf("ping mysql: %w", err) }
	return db, nil
}

package database

import (
	"example/internal/config"
	"example/pkg/logger"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	log := logger.GetLogger().With(zap.String("component", "postgres-db"))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)

	log.Info("Connecting to PostgreSQL database",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.Name),
		zap.String("user", cfg.User),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	log.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

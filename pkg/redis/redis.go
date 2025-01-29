package redis

import (
	"context"
	"example/internal/config"
	"example/pkg/logger"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	log := logger.GetLogger().With(zap.String("component", "redis-client"))

	log.Info("Connecting to Redis",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.Int("database", cfg.DB),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	log.Info("Successfully connected to Redis")
	return client, nil
}

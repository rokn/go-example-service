package cache

import (
	"context"
	"encoding/json"
	"example/pkg/logger"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redisstore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Manager defines the interface for cache operations
type Manager interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	SetDefault(ctx context.Context, key string, value interface{}) error
}

type manager struct {
	cache  *cache.Cache[any]
	logger *zap.Logger
}

func NewCacheManager(redisClient *redis.Client) Manager {
	log := logger.GetLogger().With(zap.String("component", "cache-manager"))

	redisStore := redisstore.NewRedis(redisClient,
		store.WithExpiration(1*time.Hour),
	)
	cacheManager := cache.New[any](redisStore)

	return &manager{
		cache:  cacheManager,
		logger: log,
	}
}

func (cm *manager) Get(ctx context.Context, key string, value interface{}) error {
	cm.logger.Debug("Getting value from cache", zap.String("key", key))

	result, err := cm.cache.Get(ctx, key)
	if err != nil {
		if err.Error() == store.NOT_FOUND_ERR {
			cm.logger.Debug("Cache miss", zap.String("key", key))
		} else {
			cm.logger.Error("Failed to get value from cache", zap.Error(err))
		}
		return err
	}

	// Convert result to JSON string if it's not already
	var jsonStr string
	switch v := result.(type) {
	case string:
		jsonStr = v
	default:
		bytes, err := json.Marshal(result)
		if err != nil {
			cm.logger.Error("Failed to marshal cached value", zap.Error(err))
			return err
		}
		jsonStr = string(bytes)
	}

	// Unmarshal JSON into the target value
	if err := json.Unmarshal([]byte(jsonStr), value); err != nil {
		cm.logger.Error("Failed to unmarshal cached value", zap.Error(err))
		return err
	}

	return nil
}

func (cm *manager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cm.logger.Debug("Setting value in cache",
		zap.String("key", key),
		zap.Duration("expiration", expiration),
	)

	// Marshal value to JSON string
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		cm.logger.Error("Failed to marshal value for caching", zap.Error(err))
		return err
	}

	if err := cm.cache.Set(ctx, key, string(jsonBytes), store.WithExpiration(expiration)); err != nil {
		cm.logger.Error("Failed to set value in cache", zap.Error(err))
		return err
	}

	return nil
}

func (cm *manager) Delete(ctx context.Context, key string) error {
	cm.logger.Debug("Deleting value from cache", zap.String("key", key))

	if err := cm.cache.Delete(ctx, key); err != nil {
		cm.logger.Error("Failed to delete value from cache", zap.Error(err))
		return err
	}

	return nil
}

func (cm *manager) SetDefault(ctx context.Context, key string, value interface{}) error {
	cm.logger.Debug("Setting value in cache with default expiration",
		zap.String("key", key),
	)

	// Marshal value to JSON string
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		cm.logger.Error("Failed to marshal value for caching", zap.Error(err))
		return err
	}

	if err := cm.cache.Set(ctx, key, string(jsonBytes)); err != nil {
		cm.logger.Error("Failed to set value in cache", zap.Error(err))
		return err
	}

	return nil
}

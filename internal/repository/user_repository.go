package repository

import (
	"context"
	"encoding/json"
	"example/internal/model"
	"example/pkg/logger"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		db:     db,
		redis:  redis,
		logger: logger.GetLogger().With(zap.String("component", "user-repository")),
	}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	r.logger.Info("Getting user by ID", zap.Uint("id", id))

	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d", id)

	// Check cache
	r.logger.Debug("Checking Redis cache", zap.String("key", cacheKey))
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var user model.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			r.logger.Debug("User found in cache", zap.Uint("id", id))
			return &user, nil
		}
		r.logger.Error("Failed to unmarshal cached user", zap.Error(err))
	}

	// If not in cache, get from DB
	r.logger.Debug("User not found in cache, querying database", zap.Uint("id", id))
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		r.logger.Error("Failed to get user from database", zap.Error(err))
		return nil, err
	}

	// Store in cache
	if userJSON, err := json.Marshal(user); err == nil {
		r.logger.Debug("Storing user in cache", zap.Uint("id", id))
		r.redis.Set(context.Background(), cacheKey, userJSON, 15*time.Minute)
	} else {
		r.logger.Error("Failed to marshal user for caching", zap.Error(err))
	}

	r.logger.Debug("Successfully retrieved user", zap.Uint("id", id))
	return &user, nil
}

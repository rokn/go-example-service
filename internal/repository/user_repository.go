package repository

import (
	"context"
	"errors"
	"example/internal/model"
	"example/pkg/cache"
	"example/pkg/logger"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type userRepository struct {
	db           *gorm.DB
	cacheManager cache.Manager
	logger       *zap.Logger
}

func NewUserRepository(db *gorm.DB, cacheManager cache.Manager) UserRepository {
	return &userRepository{
		db:           db,
		cacheManager: cacheManager,
		logger:       logger.GetLogger().With(zap.String("component", "user-repository")),
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	r.logger.Info("Getting user by ID", zap.Uint("id", id))

	var user model.User
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try to get from cache first
	r.logger.Debug("Checking cache", zap.String("key", cacheKey))
	err := r.cacheManager.Get(context.Background(), cacheKey, &user)
	if err == nil {
		r.logger.Debug("User found in cache", zap.Uint("id", id))
		return &user, nil
	}

	// If not in cache, get from DB
	r.logger.Debug("User not found in cache, querying database", zap.Uint("id", id))
	if err := r.db.First(&user, id).Error; err != nil {
		r.logger.Error("Failed to get user from database", zap.Error(err))
		return nil, err
	}

	// Store in cache
	r.logger.Debug("Storing user in cache", zap.Uint("id", id))
	if err := r.cacheManager.SetDefault(context.Background(), cacheKey, user); err != nil {
		r.logger.Error("Failed to cache user", zap.Error(err))
		// Don't return the error since we still have the user
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	cacheKey := fmt.Sprintf("user:email:%s", email)

	// Try to get from cache first
	err := r.cacheManager.Get(context.Background(), cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// If not in cache, get from database
	user = model.User{Email: email}
	result := r.db.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Cache the user with default expiration
	if err := r.cacheManager.SetDefault(context.Background(), cacheKey, user); err != nil {
		r.logger.Error("Failed to cache user", zap.Error(err))
		// Don't return the error since we still have the user
	}

	return &user, nil
}

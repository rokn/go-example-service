package main

import (
	"example/internal/config"
	"example/internal/http/handler"
	"example/internal/repository"
	"example/internal/router"
	"example/internal/service"
	"example/pkg/cache"
	"example/pkg/database"
	"example/pkg/logger"
	"example/pkg/redis"

	_ "example/docs" // Import swagger docs

	"go.uber.org/zap"
)

// @title           User API
// @version         1.0
// @description     A User management API with Redis caching
// @host           localhost:8080
// @BasePath       /api

func main() {
	// Initialize logger
	log, err := logger.Initialize("development")
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config", zap.Error(err))
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatal("Cannot connect to database", zap.Error(err))
	}

	// Initialize Redis
	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatal("Cannot connect to Redis", zap.Error(err))
	}

	// Initialize cache manager
	cacheManager := cache.NewCacheManager(redisClient)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db, cacheManager)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler, err := handler.NewAuthHandler(userService, &cfg.Auth)
	if err != nil {
		log.Fatal("Cannot initialize auth handler", zap.Error(err))
	}

	// Initialize and start router
	r := router.SetupRouter(userHandler, authHandler)

	// Start server
	log.Info("Starting server", zap.String("port", cfg.App.Port))
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

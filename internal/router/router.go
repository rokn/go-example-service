package router

import (
	"example/internal/http/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		panic("Failed to set trusted proxies: " + err.Error())
	}

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Middleware().LoginHandler)
			auth.GET("/refresh", authHandler.Middleware().RefreshHandler)
		}

		// Public user routes
		users := api.Group("/users")
		{
			users.POST("/", userHandler.Create) // Allow public user creation
		}

		// Protected user routes
		protected := users.Group("")
		protected.Use(authHandler.Middleware().MiddlewareFunc())
		{
			protected.GET("/:id", userHandler.Get)
		}
	}

	return r
}

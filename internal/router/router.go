package router

import (
	"example/internal/http/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		panic("Failed to set trusted proxies: " + err.Error())
	}

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", userHandler.Create) // @Summary Create user
			users.GET("/:id", userHandler.Get)  // @Summary Get user by ID
		}
	}

	return r
}

package handler

import (
	"example/internal/config"
	"example/internal/http/handler/requests"
	"example/internal/model"
	"example/internal/service"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler defines the interface for authentication handler operations
type AuthHandler interface {
	Middleware() *jwt.GinJWTMiddleware
}

type authHandler struct {
	userService    service.UserService
	authMiddleware *jwt.GinJWTMiddleware
}

func NewAuthHandler(userService service.UserService, cfg *config.AuthConfig) (AuthHandler, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           cfg.Realm,
		Key:             []byte(cfg.SecretKey),
		Timeout:         24 * time.Hour,
		MaxRefresh:      72 * time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator(userService),
		Authorizator:    authorizator,
		Unauthorized:    unauthorized,
		TokenLookup:     "header: Authorization, query: token",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		return nil, err
	}

	return &authHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}, nil
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*model.User); ok {
		return jwt.MapClaims{
			"id":    v.ID,
			"email": v.Email,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &model.User{
		Model: gorm.Model{
			ID: uint(claims["id"].(float64)),
		},
		Email: claims["email"].(string),
	}
}

func authenticator(userService service.UserService) func(*gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginReq requests.LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		user, err := userService.Login(loginReq.Email, loginReq.Password)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		if user == nil {
			return nil, jwt.ErrFailedAuthentication
		}

		return user, nil
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	return true // Authorization is now handled by route-specific middleware
}

func unauthorized(c *gin.Context, code int, message string) {
	NewErrorResponse(c, code, "Unauthorized", []interface{}{message})
}

func (h *authHandler) Middleware() *jwt.GinJWTMiddleware {
	return h.authMiddleware
}

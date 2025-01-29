package handler

import (
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
}

func NewSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, errors []interface{}) {
	c.JSON(statusCode, BaseResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})

}

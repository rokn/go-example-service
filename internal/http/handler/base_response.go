package handler

import (
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Success bool          `json:"success"`
	Status  int           `json:"status,omitempty"`
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
}

func NewSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, BaseResponse{
		Success: true,
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, errors []interface{}) {
	c.JSON(statusCode, BaseResponse{
		Success: false,
		Status:  statusCode,
		Message: message,
		Errors:  errors,
	})

}

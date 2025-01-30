package handler

import (
	"net/http"
	"strconv"

	"example/internal/http/handler/requests"
	"example/internal/http/handler/responses"
	"example/internal/model"
	"example/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler defines the interface for user handler operations
type UserHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetMe(c *gin.Context)
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body requests.UserCreateRequest true "User information"
// @Success 201 {object} BaseResponse{data=responses.UserResponse} "User created successfully"
// @Failure 400 {object} BaseResponse "Invalid request payload"
// @Failure 500 {object} BaseResponse "Internal server error"
// @Router /users [post]
func (h *userHandler) Create(c *gin.Context) {
	var req requests.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid request payload", []interface{}{err.Error()})
		return
	}

	user := req.ToModel()
	validationErrs, err := h.service.CreateUser(user)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []interface{}{err.Error()})
		return
	}
	if len(validationErrs) > 0 {
		errs := make([]interface{}, len(validationErrs))
		for i, v := range validationErrs {
			errs[i] = v
		}
		NewErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), errs)
		return
	}

	response := responses.UserResponseFromModel(user)
	NewSuccessResponse(c, http.StatusCreated, "User created successfully", response)
}

// Get godoc
// @Summary Get a user by ID
// @Description Get user details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} BaseResponse{data=responses.UserResponse} "User retrieved successfully"
// @Failure 400 {object} BaseResponse "Invalid ID"
// @Failure 404 {object} BaseResponse "User not found"
// @Router /users/{id} [get]
func (h *userHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid ID", []interface{}{err.Error()})
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "User not found", []interface{}{err.Error()})
		return
	}

	response := responses.UserResponseFromModel(user)
	NewSuccessResponse(c, http.StatusOK, "User retrieved successfully", response)
}

// GetMe godoc
// @Summary Get logged in user details
// @Description Get details of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} BaseResponse{data=responses.UserResponse} "User retrieved successfully"
// @Failure 401 {object} BaseResponse "Unauthorized"
// @Router /users/me [get]
func (h *userHandler) GetMe(c *gin.Context) {
	user, exists := c.Get(identityKey)
	if !exists {
		NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	authenticatedUser, ok := user.(*model.User)
	if !ok {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid user data", nil)
		return
	}

	userDetails, err := h.service.GetUser(authenticatedUser.ID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user details", []interface{}{err.Error()})
		return
	}

	response := responses.UserResponseFromModel(userDetails)
	NewSuccessResponse(c, http.StatusOK, "User retrieved successfully", response)
}

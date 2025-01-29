package requests

import "example/internal/model"

// UserCreateRequest represents the request payload for creating a user
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// ToModel converts UserCreateRequest to model.User
func (r *UserCreateRequest) ToModel() *model.User {
	return &model.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

package responses

import "example/internal/model"

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	CreatedAt string `json:"created_at" example:"2024-01-01 10:00:00"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01 10:00:00"`
}

// FromModel creates UserResponse from model.User
func UserResponseFromModel(user *model.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

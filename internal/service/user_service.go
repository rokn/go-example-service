package service

import (
	"example/internal/model"
	"example/internal/repository"
	"example/pkg/validator"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user *model.User) ([]validator.ValidationError, error) {
	errors := validator.ValidateStruct(user)
	if len(errors) > 0 {
		return errors, nil
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	err = s.repo.Create(user)
	return nil, err
}

func (s *UserService) GetUser(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

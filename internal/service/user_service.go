package service

import (
	"example/internal/model"
	"example/internal/repository"
	"example/pkg/validator"

	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user service operations
type UserService interface {
	CreateUser(user *model.User) ([]validator.ValidationError, error)
	GetUser(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	Login(email, password string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(user *model.User) ([]validator.ValidationError, error) {
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

func (s *userService) GetUser(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) Login(email, password string) (*model.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil
	}

	return user, nil
}

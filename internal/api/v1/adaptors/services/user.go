package services

import (
	"errors"

	"github.com/dev-rodrigobaliza/carteado/domain/model"
	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
)

type UserService struct {
	userRepository ports.IUserRepository
}

var _ ports.IUserService = (*UserService)(nil)

// NewUserService creates a new UserService
func NewUserService(userRepository ports.IUserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Add(user *request.User) (*response.User, []*response.ErrorValidation, error) {
	// basic validation
	if user.Name == "" {
		return nil, nil, errors.New("name not found")
	}
	if user.Email == "" {
		return nil, nil, errors.New("email not found")
	}
	if user.Password == "" {
		return nil, nil, errors.New("password not found")
	}
	// validate the input fields
	validationErr := validate(user)
	if validationErr != nil {
		return nil, validationErr, nil
	}
	// add user
	u := model.NewUser(user.Name, user.Email, user.Password, false)
	err := s.userRepository.Add(u)
	if err != nil {
		return nil, nil, err
	}

	us := response.User{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		IsAdmin: u.IsAdmin,
	}

	return &us, nil, nil
}

func (s *UserService) Delete() error {
	return errors.New("not implemented")
}

func (s *UserService) Get() error {
	return errors.New("not implemented")
}

func (s *UserService) Update() error {
	return errors.New("not implemented")
}

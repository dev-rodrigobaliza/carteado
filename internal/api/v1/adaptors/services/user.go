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

func (s *UserService) Get(user *request.GetUser) (*response.User, []*response.ErrorValidation, error) {
	// basic validation
	if user.ID == 0 && user.Email == "" {
		return nil, nil, errors.New("id and/or email not found")
	}
	// validate the input fields
	validationErr := validate(user)
	if validationErr != nil {
		return nil, validationErr, nil
	}
	// get user
	var u *model.User
	var err error
	if user.ID > 0 {
		u, err = s.userRepository.FindByID(user.ID)
		if err != nil || u == nil || u.ID == 0 {
			return nil, nil, errors.New("user not found")
		}
		if user.Email != "" && user.Email != u.Email {
			return nil, nil, errors.New("user not found")
		}
	} else {
		u, err = s.userRepository.FindByEmail(user.Email)
		if err != nil || u == nil || u.ID == 0 {
			return nil, nil, errors.New("user not found")
		}		
	}

	us := response.User{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		IsAdmin: u.IsAdmin,
	}

	return &us, nil, nil
}

func (s *UserService) Update() error {
	return errors.New("not implemented")
}

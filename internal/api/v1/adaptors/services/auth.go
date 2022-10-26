package services

import (
	"strconv"
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"github.com/dev-rodrigobaliza/carteado/internal/security/paseto"
)

type AuthService struct {
	authRepository ports.IAuthRepository
	userRepository ports.IUserRepository
}

var (
	Security   *paseto.PasetoMaker
	ExpireTime time.Duration
	// This line is for get feedback in case we are not implementing the interface correctly
	_ ports.IAuthService = (*AuthService)(nil)
)

// NewAuthService creates a new AuthService
func NewAuthService(authRepository ports.IAuthRepository, userRepository ports.IUserRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

// Login logs in a user
func (s *AuthService) Login(login *request.Login, ip string) (*response.Login, []*response.ErrorValidation, error) {
	// basic validation
	if login == nil {
		return nil, nil, errors.ErrInvalidLogin
	}
	if ip == "" {
		return nil, nil, errors.ErrInvalidIP
	}
	// validate the input fields
	validationErr := validate(login)
	if validationErr != nil {
		return nil, validationErr, nil
	}
	// verify if the user exists
	user, err := s.userRepository.FindByEmail(login.Email)
	if err != nil {
		return nil, nil, err
	}
	err = user.ComparePassword(login.Password)
	if err != nil {
		return nil, nil, err
	}
	// create a new access token
	accessToken, validUntil, err := s.createAccessToken(user.ID)
	if err != nil {
		return nil, nil, err
	}
	// save login
	err = s.authRepository.Login(user.ID, accessToken)
	if err != nil {
		return nil, nil, err
	}
	// make response
	u := response.User{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}
	at := response.AccessToken{
		Token:      accessToken,
		ValidUntil: *validUntil,
	}
	response := response.Login{
		User:        &u,
		AccessToken: at,
	}

	return &response, nil, nil
}

// Logout logs out a user
func (s *AuthService) Logout(userID uint64, accessToken string) error {
	// basic validation
	if userID == 0 {
		return errors.ErrInvalidUserID
	}
	if accessToken == "" {
		return errors.ErrInvalidAccessToken
	}
	// register logout and invalidate access token
	err := s.authRepository.Logout(userID, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyToken(userID, accessToken string) error {
	// basic validation
	if userID == "" {
		return errors.ErrInvalidUserID
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	if accessToken == "" {
		return errors.ErrInvalidAccessToken
	}
	// verify token
	err = s.authRepository.VerifyToken(uint64(id), accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) createAccessToken(userID uint64) (string, *time.Time, error) {
	// basic validation
	if userID == 0 {
		return "", nil, errors.ErrInvalidUserID
	}
	id := strconv.FormatUint(userID, 10)
	// create access accessToken
	accessToken, err := Security.CreateToken(id, ExpireTime)
	if err != nil {
		return "", nil, err
	}
	validUntil := time.Now().Add(ExpireTime)

	return accessToken, &validUntil, nil
}

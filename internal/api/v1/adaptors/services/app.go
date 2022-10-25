package services

import (
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/adaptors/repositories"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
)

type AppService struct {
	AuthService ports.IAuthService
	UserService ports.IUserService
}

var appService *AppService

// NewAppService creates a new AppService with all services included
func NewAppService(appRepository *repositories.AppRepository) *AppService {
	userService := NewUserService(appRepository.UserRepository)
	authService := NewAuthService(appRepository.AuthRepository, appRepository.UserRepository)

	appService = &AppService{
		AuthService: authService,
		UserService: userService,
	}

	return appService
}

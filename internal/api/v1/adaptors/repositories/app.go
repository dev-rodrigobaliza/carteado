package repositories

import (
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"gorm.io/gorm"
)

type AppRepository struct {
	// include here all repositories
	AuthRepository ports.IAuthRepository
	UserRepository ports.IUserRepository
}

// NewAppRepository returns a new instance of AppRepository with all repositories
// the connection to the database is made here
func NewAppRepository(db *gorm.DB) *AppRepository {
	authRepository := NewAuthRepository(db)
	userRepository := NewUserRepository(db)

	return &AppRepository{
		AuthRepository: authRepository,
		UserRepository: userRepository,
	}
}

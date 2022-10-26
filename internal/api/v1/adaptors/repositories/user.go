package repositories

import (
	"github.com/dev-rodrigobaliza/carteado/domain/model"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// This line is for get feedback in case we are not implementing the interface correctly
var _ ports.IUserRepository = (*UserRepository)(nil)

// NewAuthRepository creates new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Add(user *model.User) error {
	// basic validation
	if user == nil {
		return errors.ErrInvalidUser
	}
	if user.ID > 0 {
		return errors.ErrInvalidUserID
	}
	if user.Name == "" {
		return errors.ErrInvalidName
	}
	if user.Email == "" {
		return errors.ErrInvalidEmail
	}
	if user.PasswordHash == "" {
		return errors.ErrInvalidPassword
	}
	// add user
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete() error {
	return errors.ErrNotImplemented
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	// basic validation
	if email == "" {
		return nil, errors.ErrInvalidEmail
	}
	// search for user in database
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil || user.ID == 0 {
		return nil, errors.ErrNotFoundUser
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	// basic validation
	if id == 0 {
		return nil, errors.ErrInvalidUserID
	}
	// search for user in database
	var user model.User
	user.ID = id
	err := r.db.First(&user).Error
	if err != nil {
		return nil, errors.ErrNotFoundUser
	}

	return &user, nil
}

func (r *UserRepository) Update() error {
	return errors.ErrNotImplemented
}

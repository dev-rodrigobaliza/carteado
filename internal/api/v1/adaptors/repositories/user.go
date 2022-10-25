package repositories

import (
	"errors"

	"github.com/dev-rodrigobaliza/carteado/domain/model"
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
		return errors.New("user not found")
	}
	if user.ID > 0 {
		return errors.New("user id invalid")
	}
	if user.Name == "" {
		return errors.New("user name not found")
	}
	if user.Email == "" {
		return errors.New("user name not found")
	}
	if user.PasswordHash == "" {
		return errors.New("user password not found")
	}
	// add user
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete() error {
	return errors.New("not implemented")
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	// basic validation
	if email == "" {
		return nil, errors.New("email not found")
	}
	// search for user in database
	var user model.User
	r.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	// basic validation
	if id == 0 {
		return nil, errors.New("id not found")
	}
	// search for user in database
	var user model.User
	user.ID = id
	r.db.First(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *UserRepository) Update() error {
	return errors.New("not implemented")
}


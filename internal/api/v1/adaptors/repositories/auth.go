package repositories

import (
	"github.com/dev-rodrigobaliza/carteado/domain/model"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/api/v1/ports"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

// This line is for get feedback in case we are not implementing the interface correctly
var _ ports.IAuthRepository = (*AuthRepository)(nil)

// NewAuthRepository creates new instance of AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Login(userID uint64, accessToken string) error {
	// basic validation
	if userID == 0 {
		return errors.ErrInvalidUserID
	}
	if accessToken == "" {
		return errors.ErrInvalidAccessToken
	}
	// clear previous login
	var previous []model.Login
	r.db.Where("user_id = ?", userID).Find(&previous)
	for _, login := range previous {
		login.Logout = model.LogoutAutomatic
		r.db.Save(&login)
		r.db.Delete(&login)
	}
	// save login
	login := model.Login{
		UserID: userID,
		Token: accessToken,
	}
	return r.db.Save(&login).Error
}

func (r *AuthRepository) Logout(userID uint64, accessToken string) error {
	// basic validation
	if userID == 0 {
		return errors.ErrInvalidUserID
	}
	if accessToken == "" {
		return errors.ErrInvalidAccessToken
	}
	// clear logins
	var previous []model.Login
	r.db.Where("user_id = ?", userID).Find(&previous)
	// database validation
	if len(previous) == 0 {
		return errors.ErrInvalidUserID
	}

	for _, login := range previous {
		login.Logout = model.LogoutManual
		r.db.Save(&login)
		r.db.Delete(&login)
	}

	return nil
}

func (r *AuthRepository) VerifyToken(userID uint64, accessToken string) error {
	// basic validation
	if userID == 0 {
		return errors.ErrInvalidUserID
	}
	if accessToken == "" {
		return errors.ErrInvalidAccessToken
	}
	// database validation
	var login model.Login
	err := r.db.Where("user_id = ? AND token = ?", userID, accessToken).Find(&login).Error
	if err != nil {
		return err
	}
	if login.ID == 0 {
		return errors.ErrInvalidAccessToken
	}

	return nil
}
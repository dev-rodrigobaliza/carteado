package model

import (
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
)

type User struct {
	Base             `gorm:"embedded"`
	Name             string     `json:"name" gorm:"type:varchar(200);not null"`
	Email            string     `json:"email" gorm:"type:varchar(200);not null;uniqueIndex"`
	PasswordHash     string     `json:"-" gorm:"type:varchar(128);not null"`
	IsAdmin          bool       `json:"-" gorm:"default:false;not null"`
	ConfirmedEmailAt *time.Time `json:"-" gorm:"default null"`
}

// NewUser creates a new User
func NewUser(name, email, password string, isAdmin bool) *User {
	user := User{
		Name:    name,
		Email:   email,
		IsAdmin: isAdmin,
	}
	if password != "" {
		user.SetPassword(password)
	}

	return &user
}

// ComparePassword compares the provided password with the stored hash
func (u *User) ComparePassword(password string) error {
	match, err := argon2id.ComparePasswordAndHash(password, u.PasswordHash)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid password")
	}

	return nil
}

// SetPassword sets the password hash
func (u *User) SetPassword(password string) error {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)

	return nil
}

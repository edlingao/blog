package core

import (
	"errors"
	"log"

	"github.com/edlingao/internal/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin            = "admin"
	RoleUser             = "user"
	ErrorInvalidPassword = "invalid password"
	ErrorUserNotFound    = "user not found"
	ErrorPasswordEncrypt = "password encryption failed"
)

type User struct {
	ID       string `db:"id" json:"-"`
	Username string `db:"username" json:"username"`
	Password string `db:"password_hash" json:"-"`
	Role     string `db:"role" json:"role"`
}

func NewUser(username, role string) *User {
	return &User{
		Username: username,
		Role:     role,
	}
}

func (user *User) IsAdmin() bool {
	return user.Role == "admin"
}

func (user *User) NewPassword(oldPassword, newPassword string) error {
	if oldPassword != "" && !user.ValidatePassword(oldPassword) {
		return errors.New(ErrorInvalidPassword)
	}

	encryptedPassword, err := user.encryptPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = encryptedPassword
	return nil
}

func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(ErrorPasswordEncrypt + err.Error())
	}

	return string(hash), nil
}

func (user *User) GenerateToken() string {
	token, error := auth.GenerateToken(user.ID, user.Username)
	if error != nil {
		log.Fatal("Error generating token: ", error)
		return ""
	}

	return token
}

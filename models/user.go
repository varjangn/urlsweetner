package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	UUID      string `json:"uuid"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func NewUser(email, password, firstname, lastname string) (*User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	uuid := uuid.New()
	return &User{
		Id:        -1,
		Email:     email,
		Password:  string(hashedPass),
		UUID:      uuid.String(),
		FirstName: firstname,
		LastName:  lastname,
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) UpdatePassword(newPassword string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return "", err
	}
	u.Password = string(hashedPass)
	return string(hashedPass), nil
}

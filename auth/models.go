package auth

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type account struct {
	id       uuid.UUID
	name     string
	password string
	email    string
}

func newAccount(name string, password string, email string) (*account, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("error to encrypt password %v", err)
	}

	return &account{
		name:     name,
		password: string(pwd),
		email:    email,
	}, nil
}

func (a *account) comparerPassword(raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.password), []byte(raw))
	return err == nil
}

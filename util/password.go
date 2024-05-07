package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return errors.New("password is incorrect")
	default:
		return err
	}
}

package service

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func checkEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func checkString(s string) error {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return errors.New("Empty string")
	}
	return nil
}

func checkCreds(userName, email, password string) error {
	if err := checkEmail(email); err != nil {
		return err
	}
	if err := checkString(password); err != nil {
		return err
	}
	if err := checkString(userName); err != nil {
		return err
	}
	return nil
}

func generatePassHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func newToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

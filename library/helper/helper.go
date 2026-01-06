package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckHashPassword(newPass string, oldPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(oldPass), []byte(newPass))
	if err != nil {
		return err
	}
	return nil
}

package helper

import (
	"strconv"
	"strings"

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

func GenerateSlug(s string) string {
	j := strings.ToLower(s)
	slug := strings.ReplaceAll(j, " ", "-")
	return slug
}

func StringToInt(s string) (int64, error) {
	newInt, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return newInt, nil
}

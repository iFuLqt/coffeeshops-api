package helper

import (
	"fmt"
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

func LinkInstagram(s string) string {
	return fmt.Sprintf("https://instagram.com/%s", s)
}

func LinkGoogleMaps(lat, long float64) string {
	return fmt.Sprintf("https://www.google.com/maps?q=%.6f,%.6f", lat, long)
}

func GenerateOpenTime(o, c string) string {
	open := CutTime(o)
	close := CutTime(c)
	return fmt.Sprintf("%s - %s", open, close)
}

func CutTime(s string) string {
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return s
	}
	return parts[0] + ":" + parts[1]
}

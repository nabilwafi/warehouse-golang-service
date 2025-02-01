package utils

import (
	"errors"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword enkripsi kata sandi menggunakan bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), viper.GetInt("SALT"))
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hashedPassword), nil
}

// ComparePasswords membandingkan kata sandi yang diberikan dengan kata sandi yang dienkripsi
func ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

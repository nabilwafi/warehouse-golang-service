package config

import (
	"os"
	"time"
)

type DBConfig struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     string
	SSLMode  string
	Schema   string
}

type HTTPConfig struct {
	Port string
}

type Config struct {
	DB   DBConfig
	Http HTTPConfig
}

var (
	secretKey = os.Getenv("SECRET_KEY")
	expTime   = 120 * time.Minute
)

func NewEnv() (Config, error) {
	db := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		Schema:   os.Getenv("DB_SCHEMA"),
	}
	http := HTTPConfig{
		Port: os.Getenv("PORT"),
	}

	config := Config{
		DB:   db,
		Http: http,
	}

	return config, nil
}

func GetJwtSecretKey() []byte {
	return []byte(secretKey)
}

func GetExpTime() time.Duration {
	return expTime
}

package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/config"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
)

type CustomClaims struct {
	ID    uuid.UUID
	Name  string
	Email string
	Role  dto.UserRole
	jwt.RegisteredClaims
}

func GenerateToken(user *CustomClaims) (string, error) {
	exp := time.Now().Add(config.GetExpTime())

	claims := CustomClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJwtSecretKey())
}

func VerifyToken(tokenString string) (user *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return config.GetJwtSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, jwt.ErrTokenExpired
		}

		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

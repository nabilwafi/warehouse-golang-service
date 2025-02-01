package utils

import "github.com/google/uuid"

func GenerateUUID() uuid.UUID {
	uuid := uuid.New()

	return uuid
}

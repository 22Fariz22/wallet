package utils

import "github.com/google/uuid"

// ValidateUUID проверяет, является ли строка валидным UUID
func ValidateUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}

package utils

import (
	"fmt"
	"uuid"
)

// GenerateID creates a time-ordered UUIDv7
func GenerateID(prefix string) string {
	id := uuid.NewV7().String()
	return fmt.Sprintf("%s_%s", prefix, id)
}

package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateOpaqueToken(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("invalid token size")
	}

	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

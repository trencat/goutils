package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateToken generates URL-safe, base64 encoded
// securely generated random string.
func GenerateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

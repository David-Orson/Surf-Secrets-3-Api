package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken() (string, error) {
	base := 64
	bytes := make([]byte, base)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	sessionId := base64.StdEncoding.EncodeToString(bytes)

	return sessionId, nil
}

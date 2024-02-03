package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func GenerateAndHashToken(length int) (string, string, error) {
	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	token := base64.URLEncoding.EncodeToString(randomBytes)
	token = token[:length]

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), 12)

	if err != nil {
		return "", "", err
	}

	return token, string(tokenHash), nil
}

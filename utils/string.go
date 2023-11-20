package utils

import (
	"crypto/rand"
	"encoding/base64"
)

const resetTokenBytes = 32

func GenerateResetToken() (string, error) {
	tokenBytes := make([]byte, resetTokenBytes)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	resetToken := base64.URLEncoding.EncodeToString(tokenBytes)
	return resetToken, nil
}

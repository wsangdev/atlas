package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

const apiKeyBytes = 32

// GenerateAPIKey crea una API key aleatoria (64 caracteres hex).
func GenerateAPIKey() (string, error) {
	buf := make([]byte, apiKeyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// HashAPIKey guarda solo el hash: en la DB nunca se persiste la key en claro.
func HashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

// VerifyAPIKey compara en tiempo constante para evitar timing attacks.
func VerifyAPIKey(hash, key string) bool {
	expected := HashAPIKey(key)
	return subtle.ConstantTimeCompare([]byte(hash), []byte(expected)) == 1
}

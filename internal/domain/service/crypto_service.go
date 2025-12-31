package service

import "golang.org/x/crypto/bcrypt"

const apiKeyPrefixLen = 8

func HashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparAPIKey(hash, apiKey string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(apiKey)) == nil
}

func APIKeyPrefix(apiKey string) string {
	if len(apiKey) <= apiKeyPrefixLen {
		return apiKey
	}
	return apiKey[:apiKeyPrefixLen]
}

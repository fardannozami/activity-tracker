package service

import "golang.org/x/crypto/bcrypt"

func HashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparAPIKey(hash, apiKey string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(apiKey)) == nil
}

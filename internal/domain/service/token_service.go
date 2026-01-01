package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret []byte
}

func (t TokenService) Issue(clientID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": clientID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.Secret)
}

func (t TokenService) Verify(tokenStr string) (string, error) {
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return t.Secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	sub, _ := claims["sub"].(string)
	return sub, nil
}

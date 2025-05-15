package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	key []byte
}

func New(cfg Config) *JWT {
	return &JWT{cfg.GetKey()}
}

func (j *JWT) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"is_admin": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.key)
}

func (j *JWT) IsAdmin(tokenString string) (bool, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		return false, fmt.Errorf("failed in parse token: %v", err)
	}

	if !token.Valid{
		return false, fmt.Errorf("invalid token: %v", err)
	}

	ok, exist := claims["is_admin"].(bool)
	if !exist{
		return false, fmt.Errorf("token does not contain id")
	}
	return ok, nil
}
package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwtKey - секретный ключ для подписи JWT
var jwtKey = []byte("my_secret_key")

// Claims определяет структуру данных для хранения информации о пользователе в JWT
type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

// GenerateJWT генерирует JWT для указанного логина
func GenerateJWT(login string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT проверяет валидность JWT и возвращает claims, если токен валиден
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}

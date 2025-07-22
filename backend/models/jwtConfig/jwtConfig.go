// Структуры для конфигурации jwt токенов
package jwtconfig

import "github.com/golang-jwt/jwt/v5"

type AccessToken struct {
	UserID string `json:"id"`
	Email  string `json:"email"`
	Token  string `json:"token"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	UserID string `json:"id"`
	Token  string `json:"token"`
	jwt.RegisteredClaims
}

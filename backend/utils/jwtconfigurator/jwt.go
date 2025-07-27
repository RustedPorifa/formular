// Содержит управление, валидацию и создание jwt токенов
package jwtconfigurator

import (
	"context"
	godb "formular/backend/database/SQL_postgre"
	jwtconfig "formular/backend/models/jwtConfig"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Создает access токен из resresh токена
func GenerateAccessTokenFromRefresh(refreshToken string) (string, error) {

	claims, err := ValidateAccessToken(refreshToken)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user_info, dbErr := godb.GetUserInfoByID(ctx, claims.ID)
	if dbErr != nil {
		return "", dbErr
	}

	accessToken, err := GenerateAccessToken(user_info.ID, user_info.Email)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Создаёт токен access для входа
func GenerateAccessToken(userID string, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwtconfig.AccessToken{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Создаёт refresh токен для обновления access
func GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(154 * time.Hour)
	claims := &jwtconfig.RefreshToken{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateNewTokens(userID string, email string) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Проверят валиден ли access токен и возвращает его структуру
func ValidateAccessToken(tokenString string) (*jwtconfig.AccessToken, error) {
	// Создаём пустой экземпляр AccessToken для claims
	claims := &jwtconfig.AccessToken{}

	// Парсим токен напрямую в нашу структуру
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Проверяем ошибки парсинга
	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// Возвращаем распаршенные claims
	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*jwtconfig.RefreshToken, error) {
	// Создаём пустой экземпляр AccessToken для claims
	claims := &jwtconfig.RefreshToken{}

	// Парсим токен напрямую в нашу структуру
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Проверяем ошибки парсинга
	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

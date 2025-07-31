// Содержит управление, валидацию и создание jwt токенов
package jwtconfigurator

import (
	"context"
	"fmt"
	godb "formular/backend/database/SQL_postgre"
	jwtconfig "formular/backend/models/jwtConfig"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	jwtSecret = []byte(secret)
}

// Создает access токен из resresh токена
func GenerateAccessTokenFromRefresh(refreshToken string) (string, error) {

	claims, err := ValidateAccessToken(refreshToken)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user_info, dbErr := godb.GetUserInfoByID(ctx, claims.UserID)
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
	return token.SignedString(jwtSecret)
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
	return token.SignedString([]byte(jwtSecret))
}

// Создаёт пару токенов access и refresh
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
	claims := &jwtconfig.AccessToken{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	println(claims.Email, claims.UserID, claims.ExpiresAt)
	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*jwtconfig.RefreshToken, error) {
	claims := &jwtconfig.RefreshToken{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func CreateEmailVerificationToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"iss":   "my_app",
		"aud":   "email_verification",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func VerifyEmailToken(signedToken string) (string, error) {
	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("token verification failed: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", fmt.Errorf("email claim is invalid")
		}
		return email, nil
	}

	return "", fmt.Errorf("invalid token")
}

package auth

import (
	"context"
	"errors"
	"fmt"
	godb "formular/backend/database/SQL_postgre"
	user "formular/backend/models/userConfig"
	"formular/backend/utils/jwtconfigurator"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func HandleVerify(c *gin.Context) {
	access_cookies, cookieErr := c.Cookie("access_token")
	if cookieErr != nil && !errors.Is(cookieErr, http.ErrNoCookie) {
		println(cookieErr.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
		return
	} else if cookieErr == nil {
		_, jwtErr := jwtconfigurator.ValidateAccessToken(access_cookies)
		if jwtErr != nil && errors.Is(jwtErr, jwt.ErrTokenExpired) {
			refresh_cookie, cookieErr := c.Cookie("refresh_token")
			if cookieErr != nil {
				println(cookieErr.Error())
				c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
				return
			}
			new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
			if createErr != nil {
				println(createErr.Error())
				c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
				return
			}
			c.SetCookie("access_token", new_access_token, 8*60*60, "/", "127.0.0.1", false, true)
		} else if jwtErr == nil {
			println("else")
			c.JSON(http.StatusAccepted, gin.H{"verify": "true"})
			return
		}
	}
}

func HandleRegister(c *gin.Context) {
	var newUser user.User
	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("BindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	log.Printf("Received user: %+v", newUser) // Логирование полученных данных

	if newUser.Password == "" {
		log.Println("Password is empty after BindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль обязателен"})
		return
	}

	var genErr error
	newUser.ID, genErr = generateUUID()
	if genErr != nil {
		log.Printf("BindJSON error: %v", genErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID не сформирован, попробуйте ещё раз"})
		return
	}
	log.Printf("Generated user ID: %s", newUser.ID)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		log.Printf("Bcrypt error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
		return
	}

	newUser.Password = string(hashedPassword)
	newUser.IsAuthenticated = false
	log.Printf("Hashed password: %s", newUser.Password) // Логирование хеша

	newUser.Role = "Anonymous"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := godb.AddUser(ctx, &newUser); err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД"})
		return
	}

	c.Redirect(http.StatusCreated, "/api/email/verify")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь создан",
		"user_id": newUser.ID,
	})
}

func HandleLogin(c *gin.Context) {
	var credentials user.Credentials
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Для отладки
	log.Printf("Login attempt: Email=%s, Password=%s", credentials.Email, credentials.Password)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := godb.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		log.Printf("User not found: %s", credentials.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", credentials.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	tokenAccessString, err := jwtconfigurator.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	tokenRefreshString, err := jwtconfigurator.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Refresh token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// Устанавливаем куки и возвращаем токен
	c.SetCookie("refresh_token", tokenRefreshString, 60*60*24*7, "/", "127.0.0.1:8080", false, true)
	c.SetCookie("access_token", tokenAccessString, 8*60*60, "/", "127.0.0.1:8080", false, true)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
		},
	})
}

func HandleRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token отсутствует"})
		return
	}

	// Валидируем токен
	claims, err := jwtconfigurator.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный refresh token"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user_info, dbErr := godb.GetUserInfoByID(ctx, claims.ID)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка поиска пользователя в дб\n" + dbErr.Error()})
		return
	}
	accessToken, newRefreshToken, err := jwtconfigurator.GenerateNewTokens(user_info.ID, user_info.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токенов"})
		return
	}

	// Обновляем куку с refresh token
	c.SetCookie("refresh_token", newRefreshToken, 60*60*24*7, "/", "127.0.0.1:8080", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id": claims.UserID,
		},
	})
}

func HandleEmail(c *gin.Context) {

}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

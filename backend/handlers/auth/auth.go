package auth

import (
	"context"
	godb "formular/backend/database"
	user "formular/backend/models/userConfig"
	"formular/backend/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		log.Printf("Bcrypt error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
		return
	}

	newUser.Password = string(hashedPassword)
	log.Printf("Hashed password: %s", newUser.Password) // Логирование хеша

	newUser.Role = "member"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := godb.AddUser(ctx, &newUser); err != nil {
		log.Printf("Database error: %v", err) // Детальное логирование ошибки БД
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь создан"})
}
func HandleLogin(c *gin.Context) {
	println("---------------")
	users, _ := godb.GetAllUsers(context.Background())
	for _, u := range users {
		println(u.Email, u.Password)
	}
	println("----------------")
	// Получаем учетные данные
	var credentials user.Credentials
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Для отладки
	log.Printf("Login attempt: Email=%s, Password=%s", credentials.Email, credentials.Password)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Получаем пользователя по email
	user, err := godb.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		log.Printf("User not found: %s", credentials.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	// 2. Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", credentials.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	// 3. Генерируем токены
	tokenAccessString, err := utils.GenerateAccessToken(user.ID, false)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	tokenRefreshString, err := utils.GenerateRefreshToken(user.ID, false)
	if err != nil {
		log.Printf("Refresh token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// 4. Устанавливаем куки и возвращаем токен
	c.SetCookie("refresh_token", tokenRefreshString, 60*60*24*7, "/", "127.0.0.1:8080", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenAccessString})
}

package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Временное хранилище (вместо БД)
var users = []User{
	{ID: 1, Name: "John", Email: "john@test.com", Password: "pass123"},
}

func HandleRegister(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Проверка существования пользователя
	for _, u := range users {
		if u.Email == newUser.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже занят"})
			return
		}
	}

	// Добавляем пользователя
	newUser.ID = len(users) + 1
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь создан"})
}
func HandleLogin(c *gin.Context) {
	var jwtSecret = os.Getenv("JWT_SECRET")
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Поиск пользователя
	for _, user := range users {
		if user.Email == credentials.Email && user.Password == credentials.Password {
			// Создаём JWT токен
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.ID,
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})

			// Подписываем токен
			tokenString, _ := token.SignedString([]byte(jwtSecret))
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
}

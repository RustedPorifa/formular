package profile

import (
	"context"
	godb "formular/backend/database"
	"formular/backend/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleProfile(c *gin.Context) {
	// Извлекаем токен из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Проверяем формат: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	tokenString := parts[1]

	// Валидируем токен
	claims, err := utils.ValidateToken(tokenString)
	if err != nil || claims == nil {
		log.Printf("Validation error: %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userInfo, dberr := godb.GetUserInfoByID(ctx, claims.UserID)
	if dberr != nil {
		log.Printf("DB error: %v", dberr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	// Возвращаем JSON
	c.JSON(http.StatusOK, gin.H{
		"email":     userInfo.Email,
		"name":      userInfo.Name,
		"completed": userInfo.CompletedVariantsCount,
		"role":      userInfo.Role,
	})
}

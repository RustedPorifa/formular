package middleware

import (
	"fmt"
	"formular/backend/utils/jwtconfigurator"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Запрос на %s\nЗаголовки: %+v\n", c.Request.URL, c.Request.Header)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":            "Требуется заголовок Authorization",
				"received_headers": c.Request.Header, // Отправляем клиенту полученные заголовки
			})
			return
		}

		// Разделяем "Bearer" и сам токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный формат токена",
			})
			return
		}

		tokenString := parts[1]
		claims, err := jwtconfigurator.ValidateAccessToken(tokenString)
		if err != nil || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Неверный токен",
				"details": err.Error(),
			})
			return
		}

		// Добавляем информацию в контекст
		c.Set("userID", claims.UserID)

		// Продолжаем обработку
		c.Next()
	}
}

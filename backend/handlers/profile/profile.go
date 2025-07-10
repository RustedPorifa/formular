package profile

import (
	"context"
	godb "formular/backend/database"
	user "formular/backend/models/userConfig"
	"formular/backend/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleProfile(c *gin.Context) {
	var jwt user.JwtToken

	if err := c.BindJSON(&jwt); err != nil {
		log.Printf("BindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	utilsClaims, claimErr := utils.ValidateToken(jwt.Token) // Теперь поле Token
	if claimErr != nil || utilsClaims == nil {
		log.Printf("Validation error: %s", claimErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"}) // 401 уместнее
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userInfo, dberr := godb.GetUserInfoByID(ctx, utilsClaims.UserID)
	if dberr != nil {
		log.Printf("DB error: %v", dberr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	// Возвращаем JSON вместо HTML
	c.JSON(http.StatusOK, gin.H{
		"email":     userInfo.Email,
		"name":      userInfo.Name,
		"completed": userInfo.CompletedVariantsCount,
		"role":      userInfo.Role,
	})
}

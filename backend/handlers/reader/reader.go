package reader

import (
	"context"
	godb "formular/backend/database/SQL_postgre"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleGetVariant(c *gin.Context) {
	grade := c.Param("grade")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	variants, dbErr := godb.GetVariantsByClass(ctx, grade)
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"variants": variants,
	})
}

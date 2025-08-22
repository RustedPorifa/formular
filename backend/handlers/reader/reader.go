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
	println(variants)
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"variants": variants,
	})
}

func HandleVariantSearch(c *gin.Context) {
	uuid := c.Param("uuid")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	variant, dbErr := godb.GetVariantByUUID(ctx, uuid)
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}
	println("------------------------------")
	println(variant.VideoFilePath, variant.PDFFilePath)
	c.HTML(http.StatusOK, "variant.html", gin.H{
		"name":        variant.Name,
		"description": variant.Description,
		"uuid":        variant.UUID,
		"pdf_path":    variant.PDFFilePath,
		"video_path":  variant.VideoFilePath,
	})
}

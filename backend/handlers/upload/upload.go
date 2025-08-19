package upload

import (
	"context"
	"fmt"
	godb "formular/backend/database/SQL_postgre"
	variantconfig "formular/backend/models/variantConfig"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Variant struct {
	Name        string                `form:"name" binding:"required"`
	Class       string                `form:"class" binding:"required,number"`
	Subject     string                `form:"subject" binding:"required"`
	Description string                `form:"description"`
	Solved      bool                  `form:"solved"`
	PDFFile     *multipart.FileHeader `json:"-"`
	VideoFile   *multipart.FileHeader `json:"-"`
}

func UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга формы: " + err.Error()})
		return
	}
	variantCount := 0
	for i := 0; ; i++ {
		name := form.Value[fmt.Sprintf("variants[%d][name]", i)]
		if len(name) == 0 {
			break
		}
		variantCount++
	}

	variants := parseVariants(form, variantCount)
	println("Кол-во вариантов: ", len(variants), variantCount)
	for _, variant := range variants {
		uuid, err := generateUUID()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации UUID: " + err.Error()})
			return
		}
		var type_of_variant string
		println("Type of variant: ", variant.Subject)
		switch variant.Subject {
		case "алгебра":
			type_of_variant = "algebra"
		case "впр":
			type_of_variant = "vpr"
		case "огэ":
			type_of_variant = "oge"
		case "егэ":
			type_of_variant = "ege"
		case "геометрия":
			type_of_variant = "geometry"
		case "стереометрия":
			type_of_variant = "stereometry"
		default:
			log.Println("Не поддерживаемый формат")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неподдерживаемый предмет: " + variant.Subject})
			return
		}
		println(type_of_variant)
		path_to_сreate := fmt.Sprintf("frontend/templates/math/%s/%s", variant.Class, type_of_variant)
		println(path_to_сreate)
		dirCreateErr := os.MkdirAll(path_to_сreate+"/"+uuid, 0777)
		if dirCreateErr != nil {
			log.Println(dirCreateErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории: " + dirCreateErr.Error()})
			return
		}
		path_to_variant := filepath.Join(path_to_сreate, uuid)
		println(path_to_variant)
		uploadErr := c.SaveUploadedFile(variant.PDFFile, path_to_variant+"/"+variant.PDFFile.Filename)
		if uploadErr != nil {
			log.Println(uploadErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки файла: " + uploadErr.Error()})
			return
		}
		var videoFilePath string
		videoErr := c.SaveUploadedFile(variant.VideoFile, path_to_variant+"/"+variant.VideoFile.Filename)
		if videoErr != nil && variant.VideoFile.Filename == "" {
			videoFilePath = ""
		} else {
			videoFilePath = path_to_variant + "/" + variant.VideoFile.Filename
		}
		variant_config := variantconfig.VariantInfo{
			UUID:          uuid,
			Name:          variant.Name,
			Description:   variant.Description,
			Class:         variant.Class,
			Subject:       variant.Subject,
			Solved:        variant.Solved,
			PDFFilePath:   path_to_variant + "/" + variant.PDFFile.Filename,
			VideoFilePath: videoFilePath,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		dbErr := godb.AddVariant(ctx, &variant_config)
		if dbErr != nil {
			log.Println(dbErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления варианта: " + dbErr.Error()})
			return
		}
	}
	// Отправляем ответ клиенту
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Данные успешно получены",
	})
}

func parseVariants(form *multipart.Form, amount int) []Variant {
	var variants []Variant
	for i := 0; i < amount; i++ {
		variant := Variant{
			Name:        getFormString(form.Value, fmt.Sprintf("variants[%d][name]", i)),
			Class:       getFormString(form.Value, fmt.Sprintf("variants[%d][class]", i)),
			Subject:     getFormString(form.Value, fmt.Sprintf("variants[%d][subject]", i)),
			Description: getFormString(form.Value, fmt.Sprintf("variants[%d][description]", i)),
			Solved:      getFormBool(form.Value, fmt.Sprintf("variants[%d][solved]", i)),
			PDFFile:     getFormFile(form, fmt.Sprintf("variants[%d][pdf]", i)),
			VideoFile:   getFormFile(form, fmt.Sprintf("variants[%d][video]", i)),
		}

		variants = append(variants, variant)

	}
	return variants
}

// Функция для безопасного получения строки
func getFormString(form url.Values, key string) string {
	if values := form[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

// Функция для безопасного получения файла
func getFormFile(form *multipart.Form, key string) *multipart.FileHeader {
	if files := form.File[key]; len(files) > 0 {
		return files[0]
	}
	return nil
}

func getFormBool(form url.Values, key string) bool {
	value := getFormString(form, key)
	return value == "true" || value == "1" || value == "on"
}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

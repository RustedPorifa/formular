package upload

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Variant struct {
	Name        string `form:"name" binding:"required"`
	Class       int    `form:"class" binding:"required,number"`
	Subject     string `form:"subject" binding:"required"`
	Description string `form:"description"`
	Solved      bool   `form:"solved"`
}

func UploadHandler(c *gin.Context) {
	// Парсим форму вручную с увеличенным лимитом памяти
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга формы: " + err.Error()})
		return
	}

	// Выводим информацию в терминал
	fmt.Println("\n===== ДАННЫЕ ПОЛУЧЕНЫ ОТ АДМИНИСТРАТОРА =====")

	// Выводим текстовые поля
	fmt.Println("\nТекстовые данные:")
	for key, values := range form.Value {
		fmt.Printf("  - %s: %v\n", key, values)
	}

	// Выводим информацию о файлах
	fmt.Println("\nФайлы:")
	for key, files := range form.File {
		for i, file := range files {
			fmt.Printf("  - %s[%d]:\n", key, i)
			fmt.Printf("      Имя: %s\n", file.Filename)
			fmt.Printf("      Размер: %d байт\n", file.Size)
			fmt.Printf("      MIME тип: %s\n", file.Header.Get("Content-Type"))
		}
	}

	// Обрабатываем каждый вариант
	fmt.Println("\nДетали вариантов:")
	for i := 0; ; i++ {
		// Проверяем, есть ли данные для этого варианта
		name := form.Value[fmt.Sprintf("variants[%d][name]", i)]
		if len(name) == 0 {
			break // Больше нет вариантов
		}

		class := form.Value[fmt.Sprintf("variants[%d][class]", i)]
		subject := form.Value[fmt.Sprintf("variants[%d][subject]", i)]
		solved := form.Value[fmt.Sprintf("variants[%d][solved]", i)]

		fmt.Printf("\nВариант #%d:\n", i+1)
		fmt.Printf("  Название: %s\n", name[0])
		fmt.Printf("  Класс: %s\n", class[0])
		fmt.Printf("  Предмет: %s\n", subject[0])
		fmt.Printf("  Решенный: %s\n", solved[0])

		// Информация о PDF файле
		pdfFiles := form.File[fmt.Sprintf("variants[%d][pdf]", i)]
		if len(pdfFiles) > 0 {
			pdf := pdfFiles[0]
			fmt.Printf("  PDF файл: %s (%d байт)\n", pdf.Filename, pdf.Size)
			// Здесь вы можете обработать PDF файл
			// Например: c.SaveUploadedFile(pdf, fmt.Sprintf("uploads/%s", pdf.Filename))
		}

		// Информация о видео файле
		videoFiles := form.File[fmt.Sprintf("variants[%d][video]", i)]
		if len(videoFiles) > 0 {
			video := videoFiles[0]
			fmt.Printf("  Видео файл: %s (%d байт)\n", video.Filename, video.Size)
			// Здесь вы можете обработать видео файл
		}
	}

	fmt.Println("==================================\n")

	// Отправляем ответ клиенту
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Данные успешно получены",
	})
}

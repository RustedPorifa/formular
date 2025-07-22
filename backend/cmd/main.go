package main

import (
	godb "formular/backend/database"
	"formular/backend/handlers/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	errLoading := godotenv.Load("SECRETS.env")
	if errLoading != nil {
		log.Panic("Ошибка в инициализации .env: ", errLoading)
	}
	errDB := godb.InitDB()
	if errDB != nil {
		log.Panic(errDB)
	}
	router := gin.Default()

	router.LoadHTMLGlob("frontend/templates/*/*")
	router.Static("/static", "frontend/static")

	// Роуты API
	router.GET("/api/auth/check", auth.HandleAuthCheck)

	// Роуты HTML-страниц
	router.GET("/", homeHandler)
	router.GET("/login", loginHandler)
	router.GET("/test", testHandler)
	router.GET("/about", aboutHandler)
	router.GET("/contact", contactHandler)
	router.GET("/admin", adminDashboardHandler)
	router.GET("/profile", HandleHtmlProfile)

	router.Run(":8080")
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Главная страница",
		"Features": []string{
			"Подготовка к ЕГЭ",
			"Видеоуроки",
			"Практические задания",
		},
	})
}

func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func testHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "test.html", gin.H{})
}

func aboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"Title": "О проекте",
	})
}

func contactHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", gin.H{
		"Title": "Контакты",
	})
}

// Новый обработчик для админки
func adminDashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
		"Title": "Админ-панель",
	})
}

func HandleHtmlProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Title": "Профиль",
	})
}

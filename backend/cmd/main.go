package main

import (
	godb "formular/backend/database/SQL_postgre"
	"formular/backend/handlers/auth"
	"formular/backend/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var Domain string

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

	// Роуты
	router.GET("/api/verify", auth.HandleVerify)
	router.GET("/", homeHandler)
	router.GET("loginform", middleware.CSRFMiddleware(), loginHandler)
	//csrf group for post
	csrfGroup := router.Group("/api")
	csrfGroup.Use(middleware.CSRFMiddleware())

	csrfGroup.POST("/register", auth.HandleRegister)
	csrfGroup.POST("/login", auth.HandleLogin)

	authorizedGroup := router.Group("/user")
	authorizedGroup.Use(middleware.AuthMiddleware())
	authorizedGroup.GET("/profile", HandleHtmlProfile)
	router.Run(":5354")
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

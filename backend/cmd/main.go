package main

import (
	auth "formular/backend/handlers/auth"
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
	//errDB := godb.InitDB()
	//if errDB != nil {
	//	log.Panic(errDB)
	//}
	router := gin.Default()

	router.LoadHTMLGlob("frontend/templates/*/*")

	router.Static("/static", "frontend/static")

	// Роуты
	router.GET("/", homeHandler)
	router.GET("/about", aboutHandler)
	router.GET("/contacts", contactHandler)
	router.GET("/admin/dashboard", adminDashboardHandler)
	router.GET("/loginform", loginHandler)
	router.POST("/register", auth.HandleRegister)
	router.POST("/login", auth.HandleLogin)
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

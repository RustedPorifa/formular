package main

import (
	"context"
	nosqlredis "formular/backend/database/NOSQL_redis"
	godb "formular/backend/database/SQL_postgre"
	"formular/backend/handlers/auth"
	"formular/backend/handlers/cloudflare"
	"formular/backend/handlers/kims"
	"formular/backend/handlers/payments"
	"formular/backend/handlers/upload"
	"formular/backend/middleware"
	"formular/backend/utils/email"
	"formular/backend/utils/jwtconfigurator"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
)

var Domain string

func main() {
	//APIs
	errLoading := godotenv.Load("SECRETS.env")
	if errLoading != nil {
		log.Panic("Ошибка в инициализации .env: ", errLoading)
	}
	//GEO
	db, GeoErr := geoip2.Open("GeoLite.mmdb")
	if GeoErr != nil {
		log.Fatal(GeoErr)
	}
	defer db.Close()
	//INIT BLOCK
	jwtconfigurator.InitJWT()
	cloudflare.InitSecret()
	errDB := godb.InitDB()
	if errDB != nil {
		log.Panic(errDB)
	}
	email.InitEmail()
	nosqlredis.InitRedis()
	payments.InitRobokassa()
	println("ADDING ADMIN")
	adminErr := godb.AddAdmin()
	if adminErr != nil {
		log.Println(adminErr)
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		godb.DeleteAllUnauthenticatedUsers(ctx)
	}()
	//Add admin

	//BASED ROUTER
	router := gin.Default()
	//Настройка роутера
	router.MaxMultipartMemory = 10 << 30
	router.SetTrustedProxies(nil)
	//router.Use(middleware.GeoIPMiddleware(db))
	router.LoadHTMLGlob("frontend/templates/**/*.html")
	router.Static("/static", "frontend/static")

	//404
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})
	// Роуты
	router.GET("/verify", auth.HandleVerify)
	router.GET("/", homeHandler)
	router.GET("loginform", middleware.CSRFMiddleware(), loginHandler)
	router.GET("/submit", cloudflare.CloudflareHandler)
	//KIMS group
	variantsGroup := router.Group("/kims")
	variantsGroup.GET("/:grade/:type", kims.HandleGrade)
	//CSRF GROUP
	csrfGroup := router.Group("/api")
	csrfGroup.Use(middleware.CSRFMiddleware())
	csrfGroup.GET("/verify-email", verifyHandler)
	csrfGroup.POST("/verify/email", auth.HandleEmailVerify)
	csrfGroup.POST("/register", auth.HandleRegister)
	csrfGroup.POST("/login", auth.HandleLogin)
	csrfGroup.POST("/email/verify", auth.HandleEmailVerify)
	//PAYMENT
	paymentGroup := router.Group("/payment")
	paymentGroup.Use(middleware.CSRFMiddleware())
	paymentGroup.GET(":grade", payments.HandlePayment)
	paymentGroup.POST("result")
	paymentGroup.POST("success")
	paymentGroup.POST("fail")
	//AUTHORIZED ONLY
	authorizedGroup := router.Group("/user")
	authorizedGroup.Use(middleware.AuthMiddleware())
	authorizedGroup.GET("/profile", HandleHtmlProfile)
	//ADMIN ONLY
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AdminMiddleware())
	adminGroup.GET("/dashboard", adminDashboardHandler)
	adminGroup.POST("/upload-variants", upload.UploadHandler)
	router.Run(":5050")
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func verifyHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "email.html", gin.H{})
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

// Обработчик дешборд администратора
func adminDashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{})
}

func HandleHtmlProfile(c *gin.Context) {
	access_token, cookieAccessErr := c.Cookie("access_token")
	if cookieAccessErr != nil {
		c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
	}
	claims, jwtErr := jwtconfigurator.ValidateAccessToken(access_token)
	if jwtErr != nil {
		c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
	}
	println(claims.ID)
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Title": "Профиль",
		"Id":    claims.UserID,
		"Email": claims.Email,
	})
}

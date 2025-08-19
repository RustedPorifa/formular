package middleware

import (
	"errors"
	"fmt"
	csrfgenerator "formular/backend/utils/csrfGenerator"
	"formular/backend/utils/jwtconfigurator"
	"formular/backend/utils/tokenchecker"
	"log"
	"net"
	"net/http"
	"net/url"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oschwald/geoip2-golang"
)

const (
	csrfCookieName = "XSRF-TOKEN"
	csrfHeaderName = "X-XSRF-TOKEN"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		access_cookies, cookieErr := c.Cookie("access_token")
		if cookieErr != nil && errors.Is(cookieErr, http.ErrNoCookie) {
			refresh_cookie, cookieErr := c.Cookie("refresh_token")
			if cookieErr != nil {
				c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
				return
			}
			new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
			if createErr != nil {
				c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
				return
			}
			c.SetCookie("access_token", new_access_token, 8*60*60, "/", "127.0.0.1", false, true)
			c.Next()
		} else if cookieErr != nil && !errors.Is(cookieErr, http.ErrNoCookie) {
			log.Println(cookieErr)
			c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
			return
		} else {
			_, jwtErr := jwtconfigurator.ValidateAccessToken(access_cookies)
			if jwtErr != nil && errors.Is(jwtErr, jwt.ErrTokenExpired) {
				refresh_cookie, cookieErr := c.Cookie("refresh_token")
				if cookieErr != nil {
					c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
					return
				}
				new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
				if createErr != nil {
					log.Println(createErr)
					c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
					return
				}
				c.SetCookie("access_token", new_access_token, 8*60*60, "/", "127.0.0.1", false, true)
				c.Next()
			} else if jwtErr == nil {
				c.Next()
			} else {
				c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
				return
			}
		}
	}
}

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Для безопасных методов (GET, HEAD, OPTIONS) устанавливаем токен в куки
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			println("GET запросик")
			cookieInfo, csrfCookieErr := c.Cookie(csrfCookieName)
			println(cookieInfo)
			if csrfCookieErr == nil {
				c.Next()
				return
			} else {
				println(csrfCookieErr.Error())
				token, err := csrfgenerator.GenerateCSRFToken()
				if err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie(
					csrfCookieName,
					token,
					3600, // Время жизни
					"/",
					"",
					false, // Secure (только HTTPS)
					false, // HttpOnly=false (чтобы JS мог прочитать)
				)
				c.Next()
				return
			}

		}
		println("Not good method")
		headerToken := c.GetHeader(csrfHeaderName)
		cookieToken, errCokie := c.Cookie(csrfCookieName)

		// Декодируем URL-кодированный заголовок
		decodedHeaderToken, err := url.QueryUnescape(headerToken)
		if err != nil {
			log.Printf("CSRF header decode error: %v", err)
			decodedHeaderToken = headerToken // Используем как есть
		}

		log.Printf("Decoded header: %s, Cookie: %s", decodedHeaderToken, cookieToken)

		if errCokie != nil || decodedHeaderToken == "" || decodedHeaderToken != cookieToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			return
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_info_db, usErr := tokenchecker.ValidateAccessTokenWithRefresh(c)
		if usErr != nil {
			log.Println(usErr)
			c.Abort()
			c.HTML(http.StatusUnauthorized, "403.html", gin.H{})
			return
		}
		if user_info_db.Role == "Admin" {
			c.Next()
		} else {
			c.HTML(http.StatusUnauthorized, "403.html", gin.H{})
			c.Abort()
			return
		}
	}
}

func variantsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		grade := c.Param("grade")
		isSolved := c.Param("isSolved")
		if isSolved == "solved" {
			user_info, checkErr := tokenchecker.ValidateAccessTokenWithRefresh(c)
			if checkErr != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}
			if slices.Contains(user_info.PurchasedGrades, grade) {
				c.Next()
			} else {
				c.HTML(http.StatusForbidden, "403_pay.html", gin.H{})
				c.Abort()
				return
			}
		}

	}
}

func GeoIPMiddleware(db *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := net.ParseIP(c.ClientIP())
		fmt.Println("Client IP:", clientIP)
		if clientIP == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid IP"})
			return
		}

		record, err := db.Country(clientIP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "GeoIP lookup error"})
			return
		}
		if record == nil || record.Country.IsoCode == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Country IsoCode not found"})
			return
		}
		fmt.Println("Client country:", record.Country.IsoCode)
		if record.Country.IsoCode != "RU" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.Next()
	}
}

package middleware

import (
	"errors"
	csrfgenerator "formular/backend/utils/csrfGenerator"
	"formular/backend/utils/jwtconfigurator"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "не найден refresh токен",
				})
				return
			}
			new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
			if createErr != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Ошибка создания токена\n" + createErr.Error(),
				})
			}
			c.SetCookie("access_token", new_access_token, 8*60*60, "/", "127.0.0.1", false, true)
			c.Next()
		} else if cookieErr != nil && !errors.Is(cookieErr, http.ErrNoCookie) {
			log.Println(cookieErr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Ошибка нахождения куки файлов, зарегестрируйтесь",
			})
			return
		} else {
			_, jwtErr := jwtconfigurator.ValidateAccessToken(access_cookies)
			if jwtErr != nil && errors.Is(jwtErr, jwt.ErrTokenExpired) {
				refresh_cookie, cookieErr := c.Cookie("refresh_token")
				if cookieErr != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "не найден refresh токен",
					})
					return
				}
				new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
				if createErr != nil {
					log.Println(createErr)
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Ошибка создания токена: " + createErr.Error(),
					})
					return
				}
				c.SetCookie("access_token", new_access_token, 8*60*60, "/", "127.0.0.1", false, true)
				c.Next()
			} else if jwtErr == nil {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Время авторизации истекло",
				})
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
				token, err := csrfgenerator.GenerateCSRFToken()
				if err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				c.SetCookie(
					csrfCookieName,
					token,
					3600, // Время жизни
					"/",
					"",    // Домен (оставьте пустым для текущего домена)
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

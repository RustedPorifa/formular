package auth

import (
	"context"
	"errors"
	"fmt"
	nosqlredis "formular/backend/database/NOSQL_redis"
	godb "formular/backend/database/SQL_postgre"
	user "formular/backend/models/userConfig"
	"formular/backend/utils/email"
	"formular/backend/utils/jwtconfigurator"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type EmailVerify struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func HandleVerify(c *gin.Context) {
	access_cookies, cookieErr := c.Cookie("access_token")
	if cookieErr != nil && !errors.Is(cookieErr, http.ErrNoCookie) {
		println(cookieErr.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
		return
	} else if cookieErr == nil {
		_, jwtErr := jwtconfigurator.ValidateAccessToken(access_cookies)
		if jwtErr != nil && errors.Is(jwtErr, jwt.ErrTokenExpired) {
			refresh_cookie, cookieErr := c.Cookie("refresh_token")
			if cookieErr != nil {
				println(cookieErr.Error())
				c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
				return
			}
			new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
			if createErr != nil {
				println(createErr.Error())
				c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
				return
			}
			c.SetCookie("access_token", new_access_token, 8*60*60, "/", "formulyarka.ru", true, true)
		} else if jwtErr == nil {
			println("else")
			c.JSON(http.StatusAccepted, gin.H{"verify": "true"})
			return
		}
	} else if errors.Is(cookieErr, http.ErrNoCookie) {
		println("NO COOKIE")
		refresh_cookie, refreshErr := c.Cookie("refresh_token")
		if refreshErr != nil {
			println(refreshErr.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
			return
		}
		new_access_token, createErr := jwtconfigurator.GenerateAccessTokenFromRefresh(refresh_cookie)
		if createErr != nil {
			println("GENERATE ERROR: ", createErr.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"verify": "false"})
			return
		}
		c.SetCookie("access_token", new_access_token, 8*60*60, "/", "formulyarka.ru", true, true)
		c.JSON(http.StatusAccepted, gin.H{"verify": "true"})
		return
	}
}

func HandleRegister(c *gin.Context) {
	var newUser user.User

	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("BindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	log.Printf("Received user: %+v", newUser) // Логирование полученных данных

	if newUser.Password == "" {
		log.Println("Password is empty after BindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль обязателен"})
		return
	}

	var genErr error
	newUser.ID, genErr = generateUUID()
	if genErr != nil {
		log.Printf("BindJSON error: %v", genErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID не сформирован, попробуйте ещё раз"})
		return
	}
	log.Printf("Generated user ID: %s", newUser.ID)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		log.Printf("Bcrypt error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
		return
	}

	newUser.Password = string(hashedPassword)
	newUser.IsAuthenticated = false
	log.Printf("Hashed password: %s", newUser.Password) // Логирование хеша
	email_lower := strings.ToLower(newUser.Email)

	code := generateVerificationCode()

	redisErr := nosqlredis.Redidb.Set(c, "verification:code:"+email_lower, code, 10*time.Minute).Err()
	if redisErr != nil {
		c.JSON(500, gin.H{"error": "Ошибка Redis"})
	}

	newUser.Role = "Anonymous"
	newUser.IsAuthenticated = false
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := godb.AddUser(ctx, &newUser); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "Такой пользователь уже существует!"})
			return
		} else {
			log.Printf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД"})
			return
		}

	}
	go email.SendEmailToVerify(newUser.Email, code)
	email_token, jwtErr := jwtconfigurator.CreateEmailVerificationToken(email_lower)
	if jwtErr != nil {
		log.Printf("JWT create err: %v", jwtErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Регистрация успешна",
		"email":    email_lower,
		"redirect": "/api/verify-email?email=" + email_token,
	})
}

func HandleLogin(c *gin.Context) {
	var credentials user.Credentials
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	log.Printf("Login attempt: Email=%s, Password=%s", credentials.Email, credentials.Password)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := godb.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		log.Printf("User not found: %s", credentials.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", credentials.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}

	tokenAccessString, err := jwtconfigurator.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	tokenRefreshString, err := jwtconfigurator.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Refresh token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// Устанавливаем куки и возвращаем токен
	c.SetCookie("refresh_token", tokenRefreshString, 60*60*24*7, "/", "formulyarka.ru", true, true)
	c.SetCookie("access_token", tokenAccessString, 8*60*60, "/", "formulyarka", true, true)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
		},
	})
}

func HandleRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token отсутствует"})
		return
	}

	// Валидируем токен
	claims, err := jwtconfigurator.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный refresh token"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user_info, dbErr := godb.GetUserInfoByID(ctx, claims.ID)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка поиска пользователя в дб\n" + dbErr.Error()})
		return
	}
	accessToken, newRefreshToken, err := jwtconfigurator.GenerateNewTokens(user_info.ID, user_info.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токенов"})
		return
	}

	// Обновляем куку с refresh token
	c.SetCookie("refresh_token", newRefreshToken, 60*60*24*7, "/", "formulyarka.ru", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id": claims.UserID,
		},
	})
}

func HandleEmailVerify(c *gin.Context) {
	var VerifyInfo EmailVerify
	if err := c.BindJSON(&VerifyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	user_email, parseErr := jwtconfigurator.VerifyEmailToken(VerifyInfo.Email)
	if parseErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка верификации токена: " + parseErr.Error()})
		return
	}
	println(VerifyInfo.Email, VerifyInfo.Code, user_email)
	redi_code, rediErr := nosqlredis.GetVerificationCode(user_email, c)
	if rediErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка датабазы redis: " + rediErr.Error()})
		return
	}
	println(redi_code)
	if redi_code == VerifyInfo.Code {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		finded_user, dbErr := godb.GetUserByEmail(ctx, user_email)
		if dbErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка датабазы: " + dbErr.Error()})
			return
		}
		godb.SetUserAuthenticatedAndRole(ctx, finded_user.ID)
		println(finded_user.ID)
		tokenAccessString, err := jwtconfigurator.GenerateAccessToken(finded_user.ID, finded_user.Email)
		if err != nil {
			log.Printf("Token generation error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
			return
		}

		tokenRefreshString, err := jwtconfigurator.GenerateRefreshToken(finded_user.ID)
		if err != nil {
			log.Printf("Refresh token generation error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
			return
		}
		c.SetCookie("refresh_token", tokenRefreshString, 60*60*24*7, "/", "formulyarka.ru", true, true)
		c.SetCookie("access_token", tokenAccessString, 8*60*60, "/", "formulyarka.ru", true, true)
		c.JSON(http.StatusCreated, gin.H{
			"message":  "Регистрация успешна! Добро пожаловать!",
			"redirect": "/",
		})
	} else {
		c.JSON(400, gin.H{"error": "Неверный код авторизации"})
		return
	}

}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

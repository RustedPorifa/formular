package nosqlredis

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var Redidb *redis.Client

func InitRedis() {
	Redidb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis
		Password: "",               // Пароль (если есть)
		DB:       0,                // Номер базы
	})

	// Проверка подключения
	_, err := Redidb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
}

func GetVerificationCode(email string, ctx *gin.Context) (string, error) {
	email_lower := strings.ToLower(email)

	code, err := Redidb.Get(ctx, "verification:code:"+email_lower).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("код не найден или истек")
	} else if err != nil {
		return "", fmt.Errorf("ошибка Redis: %v", err)
	}

	return code, nil
}

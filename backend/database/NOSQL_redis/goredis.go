package nosqlredis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis
		Password: "",               // Пароль (если есть)
		DB:       0,                // Номер базы
	})

	// Проверка подключения
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
}

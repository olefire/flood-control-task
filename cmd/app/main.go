package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"task/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if redisStatus := redisClient.Ping(ctx); redisStatus.Err() != nil {
		log.Fatal("No Redis: %w", redisStatus.Err())
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

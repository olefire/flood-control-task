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

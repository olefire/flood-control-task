package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"task/internal/config"
	"task/internal/ratelimiter"
	"time"
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

	floodLimiter := ratelimiter.New(
		ratelimiter.Deps{
			RedisClient: redisClient,
		},
		ratelimiter.Opts{
			Burst:  cfg.Burst,
			Rate:   cfg.Rate,
			Window: cfg.Window,
		},
	)

	for i := 0; i < 50; i++ {
		if result, err := floodLimiter.Check(ctx, 1); err != nil {
			log.Printf("Error while performing rate limit check: %v\n", err)
		} else {
			if result == true {
				log.Println("Check passed")
			} else {
				log.Println("Check failed")
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

}

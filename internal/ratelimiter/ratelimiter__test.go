package ratelimiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"task/internal/config"
	"testing"
	"time"
)

func TestLimiterClient_Check(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if redisStatus := redisClient.Ping(ctx); redisStatus.Err() != nil {
		t.Error(redisStatus.Err())
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			t.Error(err)
		}
	}()

	limiterClient := New(
		Deps{
			RedisClient: redisClient,
		},
		Opts{
			Burst:  cfg.Burst,
			Rate:   cfg.Rate,
			Window: cfg.Window,
		})

	for i := 0; i < 50; i++ {
		if result, err := limiterClient.Check(ctx, 1); err != nil {
			t.Logf("Error while performing rate limit check: %v\n", err)
		} else {
			if result == true {
				t.Log("Check passed")
			} else {
				t.Log("Check failed")
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

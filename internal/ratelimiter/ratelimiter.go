package ratelimiter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ FloodControl = (*FloodLimiter)(nil)

type Deps struct {
	RedisClient *redis.Client
}

type Opts struct {
	Burst  int64
	Rate   float64
	Window time.Duration
}

type FloodLimiter struct {
	deps Deps
	opts Opts
}

func New(d Deps, o Opts) *FloodLimiter {
	return &FloodLimiter{
		deps: d,
		opts: o,
	}
}

type bucket struct {
	Tokens     int64     `json:"tokens"`
	LastAccess time.Time `json:"lastAccess"`
}

func (l FloodLimiter) Check(ctx context.Context, userID int64) (bool, error) {
	bucketKey := fmt.Sprintf("user-%d", userID)

	allowed := false
	rTx := func(tx *redis.Tx) error {
		redisTime, err := tx.Time(ctx).Result()
		if err != nil {
			return fmt.Errorf("failed to get redis local time: %w", err)
		}

		var b bucket
		if rawBucket, err := tx.Get(ctx, bucketKey).Result(); errors.Is(err, redis.Nil) {
			b.Tokens = l.opts.Burst
			b.LastAccess = redisTime
		} else if err != nil {
			return fmt.Errorf("failed to receive bucket: %w", err)
		} else {
			if err = json.Unmarshal([]byte(rawBucket), &b); err != nil {
				return fmt.Errorf("failed to unmarshal bucket: %w", err)
			}

			var elapsed float64
			if b.LastAccess.Before(redisTime) {
				elapsed = float64(redisTime.Sub(b.LastAccess) / time.Second)
			}

			addTokens := int64(elapsed * l.opts.Rate / l.opts.Window.Seconds())
			remainsToNextAccess := time.Duration(math.Ceil(float64(addTokens)/l.opts.Rate)) * time.Second

			b.Tokens = min(b.Tokens+addTokens, l.opts.Burst)
			b.LastAccess = b.LastAccess.Add(remainsToNextAccess)
		}

		if b.Tokens > 0 {
			allowed = true
			b.Tokens -= 1
		}

		payload, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("failed to marshal new bucket: %w", err)
		}

		return tx.SetEx(ctx, bucketKey, payload, l.opts.Window).Err()
	}
	if err := l.deps.RedisClient.Watch(ctx, rTx, fmt.Sprint(userID)); err != nil {
		return false, fmt.Errorf("failed to run redis tx: %w", err)
	}

	return allowed, nil
}

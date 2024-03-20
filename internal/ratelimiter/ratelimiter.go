package ratelimiter

import (
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	*redis.Client
	FloodControl
}

type LimiterClient struct {
	Deps
}

func New(d Deps) *LimiterClient {
	return &LimiterClient{
		Deps: d,
	}
}

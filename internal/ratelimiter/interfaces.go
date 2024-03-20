package ratelimiter

import "context"

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

type Redis interface{}

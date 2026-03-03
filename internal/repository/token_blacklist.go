package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Blacklist struct {
	rdb *redis.Client
}

func NewTokenBlacklist(rdb *redis.Client) *Blacklist {
	return &Blacklist{rdb: rdb}
}

func (r *Blacklist) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	err := r.rdb.Set(ctx, "blacklist:"+jti, 1, ttl).Err()
	if err != nil {
		return fmt.Errorf("Revoke: %w", err)
	}
	return nil
}

func (r *Blacklist) IsRevoked(ctx context.Context, jti string) (bool, error) {
	err := r.rdb.Get(ctx, "blacklist:"+jti).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("IsRevoked: %w", err)
	}
	return true, nil
}

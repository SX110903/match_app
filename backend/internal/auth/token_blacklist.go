package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ITokenBlacklist interface {
	Add(ctx context.Context, jti string, expiry time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

type redisTokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist(client *redis.Client) ITokenBlacklist {
	return &redisTokenBlacklist{client: client}
}

func (b *redisTokenBlacklist) Add(ctx context.Context, jti string, expiry time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", jti)
	if err := b.client.Set(ctx, key, "1", expiry).Err(); err != nil {
		return fmt.Errorf("adding token to blacklist: %w", err)
	}
	return nil
}

func (b *redisTokenBlacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", jti)
	result, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("checking token blacklist: %w", err)
	}
	return result > 0, nil
}

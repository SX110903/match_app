package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type RateLimiter struct {
	redis    *redis.Client
	limit    int
	window   time.Duration
	keyFunc  func(r *http.Request) string
}

// NewIPRateLimiter creates a rate limiter keyed by client IP.
func NewIPRateLimiter(redisClient *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := &RateLimiter{
		redis:  redisClient,
		limit:  limit,
		window: window,
		keyFunc: func(r *http.Request) string {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			return fmt.Sprintf("rl:ip:%s", ip)
		},
	}
	return rl.middleware
}

// NewEndpointRateLimiter creates a rate limiter for a specific endpoint + IP combination.
func NewEndpointRateLimiter(redisClient *redis.Client, endpoint string, limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := &RateLimiter{
		redis:  redisClient,
		limit:  limit,
		window: window,
		keyFunc: func(r *http.Request) string {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			return fmt.Sprintf("rl:%s:%s", endpoint, ip)
		},
	}
	return rl.middleware
}

func (rl *RateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.keyFunc(r)
		ctx := context.Background()

		count, err := rl.redis.Incr(ctx, key).Result()
		if err != nil {
			// Fail open: allow request if Redis is down
			next.ServeHTTP(w, r)
			return
		}

		if count == 1 {
			rl.redis.Expire(ctx, key, rl.window)
		}

		if count > int64(rl.limit) {
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", rl.window.Seconds()))
			response.TooManyRequests(w, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

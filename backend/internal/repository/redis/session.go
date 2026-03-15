package redis

import (
	"context"
	"time"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type sessionRepo struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) domain.SessionRepository {
	return &sessionRepo{client: client}
}

func (r *sessionRepo) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *sessionRepo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

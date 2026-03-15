package domain

import (
	"context"
	"time"
)

type SessionContext struct {
	SessionID string `json:"sessionId"`
	Email     string `json:"email"`
}

type SessionRepository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

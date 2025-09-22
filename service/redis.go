package service

import (
	"context"
	"time"

	"speakpall/pkg/logger"
	"speakpall/storage"
)

const refreshTokenPrefix = "refresh_token:"

type RedisService interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}

type redisService struct {
	redis storage.IRedisStorage
	log   logger.ILogger
}

func NewRedisService(redis storage.IRedisStorage, log logger.ILogger) RedisService {
	return &redisService{redis: redis, log: log}
}

func (s *redisService) SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	// ixtiyoriy log — juda shovqin qilmasa foydali
	s.log.Info("redis SetX", logger.String("key", key))
	// Agar yuqoridan deadline kelmagan bo'lsa, kichik timeout beramiz
	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}
	return s.redis.SetX(ctx, key, value, duration)
}

func (s *redisService) Get(ctx context.Context, key string) (string, error) {
	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}
	return s.redis.Get(ctx, key)
}

func (s *redisService) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := refreshTokenPrefix + userID
	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}
	return s.redis.Delete(ctx, key)
}

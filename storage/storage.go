package storage

import (
	"context"
	"time"

	"speakpall/api/models"
)

type IStorage interface {
	User() IUserStorage

	Redis() IRedisStorage
	Close()
}

type IUserStorage interface {
	Create(ctx context.Context, req models.SignupRequest) (string, error)

	GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error)

	GetByID(ctx context.Context, id string) (*models.User, error)

	UpdatePassword(ctx context.Context, userID, newPassword string) error

	GetPasswordByID(ctx context.Context, userID string) (string, error)






}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error // ⬅️ YANGI
}

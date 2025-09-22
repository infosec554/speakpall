// storage/interfaces.go
package storage

import (
	"context"
	"time"

	"speakpall/api/models"
)

type IStorage interface {
	User() IUserStorage
	Profile() IProfileStorage
	Redis() IRedisStorage
	Close()
}

type IUserStorage interface {
	CreateUser(ctx context.Context, req models.SignupRequest) (string, error)
	GetLoginByEmail(ctx context.Context, email string) (models.LoginUser, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdatePasswordHash(ctx context.Context, userID, newHash string) error
	UpdateRole(ctx context.Context, userID, role string) error
	GetPasswordByID(ctx context.Context, userID string) (string, error)

	SavePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	GetUserIDByPasswordResetToken(ctx context.Context, token string, now time.Time) (string, error)
}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}


type IProfileStorage interface {
	GetProfile(ctx context.Context, userID string) (*models.Profile, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error

	GetUserInterests(ctx context.Context, userID string) ([]int, error)
	ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error

	GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error)
	UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error

	GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error)
	UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error
}
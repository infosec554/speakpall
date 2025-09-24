package service

import (
	"context"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/storage"
)

type SettingsService interface {
	GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error)
	UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error
}

type settingsService struct {
	stg storage.ISettingsStorage
	log logger.ILogger
}

func NewSettingsService(stg storage.IStorage, log logger.ILogger) SettingsService {
	return &settingsService{
		stg: stg.Settings(),
		log: log,
	}
}

func (s *settingsService) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	s.log.Info("SettingsService.GetUserSettings", logger.String("user_id", userID))
	return s.stg.GetUserSettings(ctx, userID)
}

func (s *settingsService) UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error {
	s.log.Info("SettingsService.UpsertUserSettings", logger.String("user_id", userID))
	// business rules can be applied here:
	// e.g. if req.NotifyEmail == true then ensure user's email_verified (requires User storage lookup)
	return s.stg.UpsertUserSettings(ctx, userID, req)
}

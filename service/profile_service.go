package service

import (
	"context"
	"fmt"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/pkg/profileutil"
	"speakpall/storage"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*models.Profile, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error


}

type profileService struct {
	stg storage.IProfileStorage
	log logger.ILogger
}

func NewProfileService(stg storage.IStorage, log logger.ILogger) ProfileService {
	return &profileService{
		stg: stg.Profile(),
		log: log,
	}
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (*models.Profile, error) {
	s.log.Info("ProfileService.GetProfile", logger.String("user_id", userID))
	return s.stg.GetProfile(ctx, userID)
}

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error {
	s.log.Info("ProfileService.UpdateProfile", logger.String("user_id", userID))
	patch := req

	// normalize/validate via profileutil
	if patch.DisplayName != nil {
		n, err := profileutil.NormalizeDisplayName(*patch.DisplayName)
		if err != nil {
			return fmt.Errorf("invalid display name: %w", err)
		}
		patch.DisplayName = &n
	}
	if patch.CountryCode != nil {
		cc, err := profileutil.NormalizeCountryCode(*patch.CountryCode)
		if err != nil {
			return fmt.Errorf("invalid country code: %w", err)
		}
		patch.CountryCode = &cc
	}
	if patch.Gender != nil {
		g, err := profileutil.NormalizeGender(*patch.Gender)
		if err != nil {
			return fmt.Errorf("invalid gender: %w", err)
		}
		patch.Gender = &g
	}
	if err := profileutil.ValidateAgePtr(patch.Age); err != nil {
		return err
	}
	if err := profileutil.ValidateLevelPtr(patch.Level); err != nil {
		return err
	}

	return s.stg.UpdateProfile(ctx, userID, patch)
}

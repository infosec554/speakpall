package service

import (
	"context"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/pkg/profileutil"
	"speakpall/storage"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*models.Profile, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error

	GetUserInterests(ctx context.Context, userID string) ([]int, error)
	ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error

	GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error)
	UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error

	GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error)
	UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error
}

type profileService struct {
	stg storage.IProfileStorage
	log logger.ILogger
}

func NewProfileService(stg storage.IStorage, log logger.ILogger) ProfileService {
	return &profileService{stg: stg.Profile(), log: log}
}

// -------- PROFILE --------

func (s *profileService) GetProfile(ctx context.Context, userID string) (*models.Profile, error) {
	s.log.Info("ProfileService.GetProfile", logger.String("user_id", userID))
	return s.stg.GetProfile(ctx, userID)
}

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error {
	s.log.Info("ProfileService.UpdateProfile", logger.String("user_id", userID))
	patch := req

	// normalize/validate via pkg
	if patch.DisplayName != nil {
		n, err := profileutil.NormalizeDisplayName(*patch.DisplayName)
		if err != nil { return err }
		patch.DisplayName = &n
	}
	if patch.CountryCode != nil {
		cc, err := profileutil.NormalizeCountryCode(*patch.CountryCode)
		if err != nil { return err }
		patch.CountryCode = &cc
	}
	if patch.Gender != nil {
		g, err := profileutil.NormalizeGender(*patch.Gender)
		if err != nil { return err }
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

// -------- INTERESTS --------

func (s *profileService) GetUserInterests(ctx context.Context, userID string) ([]int, error) {
	s.log.Info("ProfileService.GetUserInterests", logger.String("user_id", userID))
	return s.stg.GetUserInterests(ctx, userID)
}

func (s *profileService) ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error {
	s.log.Info("ProfileService.ReplaceUserInterests",
		logger.String("user_id", userID), logger.Int("count", len(interestIDs)),
	)
	ids := profileutil.DedupInts(interestIDs)
	return s.stg.ReplaceUserInterests(ctx, userID, ids)
}

// -------- SETTINGS --------

func (s *profileService) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	s.log.Info("ProfileService.GetUserSettings", logger.String("user_id", userID))
	return s.stg.GetUserSettings(ctx, userID)
}

func (s *profileService) UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error {
	s.log.Info("ProfileService.UpsertUserSettings", logger.String("user_id", userID))
	// Biznes qoidalari bo'lsa shu yerda (masalan email verify bo'lsa notify_email= true)
	return s.stg.UpsertUserSettings(ctx, userID, req)
}

// -------- MATCH PREFS --------

func (s *profileService) GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error) {
	s.log.Info("ProfileService.GetMatchPrefs", logger.String("user_id", userID))
	return s.stg.GetMatchPrefs(ctx, userID)
}

func (s *profileService) UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error {
	s.log.Info("ProfileService.UpsertMatchPrefs", logger.String("user_id", userID))
	patch := req

	// 1) min/max daraja tekshiruvi
	if err := profileutil.ValidateLevelRange(patch.MinLevel, patch.MaxLevel); err != nil {
		return err
	}

	// 2) gender_filter normalize
	if patch.GenderFilter != nil {
		gf, err := profileutil.NormalizeGenderFilter(*patch.GenderFilter)
		if err != nil {
			return err
		}
		patch.GenderFilter = &gf
	}

	// 3) countries_allow normalize (E'TIBOR: bu slice, pointer emas!)
	if patch.CountriesAllow != nil { // nil => key umuman yuborilmagan
		ccs := profileutil.NormalizeCountries(patch.CountriesAllow)
		patch.CountriesAllow = ccs
	}

	return s.stg.UpsertMatchPrefs(ctx, userID, patch)
}

package service

import (
	"context"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/pkg/profileutil"
	"speakpall/storage"
)

type MatchsService interface {
	GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error)
	UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error
}

type matchsService struct {
	stg storage.IMatchPreferencesStorage
	log logger.ILogger
}

func NewMatchsService(stg storage.IStorage, log logger.ILogger) MatchsService {
	return &matchsService{
		stg: stg.Matchs(),
		log: log,
	}
}


func (s *matchsService) GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error) {
	s.log.Info("InteresService.GetMatchPrefs", logger.String("user_id", userID))
	return s.stg.GetMatchPrefs(ctx, userID)
}

func (s *matchsService) UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error {
	s.log.Info("ItresService.UpsertMatchPrefs", logger.String("user_id", userID))
	if err := profileutil.ValidateLevelRange(req.MinLevel, req.MaxLevel); err != nil {
		return err
	}

	if req.GenderFilter != nil {
		gf, err := profileutil.NormalizeGenderFilter(*req.GenderFilter)
		if err != nil {
			return err
		}
		req.GenderFilter = &gf
	}

	if req.CountriesAllow != nil {
		req.CountriesAllow = profileutil.NormalizeCountries(req.CountriesAllow)
	}

	return s.stg.UpsertMatchPrefs(ctx, userID, req)
}

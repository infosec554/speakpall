package service

import (
	"context"

	"speakpall/pkg/logger"
	"speakpall/pkg/profileutil"
	"speakpall/storage"
)

type InteresService interface {
	GetUserInterests(ctx context.Context, userID string) ([]int, error)
	ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error
}

type interesService struct {
	stg storage.IUserInterestsStorage
	log logger.ILogger
}

func NewInteresService(stg storage.IStorage, log logger.ILogger) InteresService {
	return &interesService{
		stg: stg.Interest(),
		log: log,
	}
}

func (s *interesService) GetUserInterests(ctx context.Context, userID string) ([]int, error) {
	s.log.Info("ProfileService.GetUserInterests", logger.String("user_id", userID))
	return s.stg.GetUserInterests(ctx, userID)
}

func (s *interesService) ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error {
	s.log.Info("ProfileService.ReplaceUserInterests",
		logger.String("user_id", userID), logger.Int("count", len(interestIDs)),
	)
	ids := profileutil.DedupInts(interestIDs)
	return s.stg.ReplaceUserInterests(ctx, userID, ids)
}

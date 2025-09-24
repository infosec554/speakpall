package service

import (
	"context"
	"fmt"

	"speakpall/pkg/logger"
	"speakpall/storage"
)

type FriendService interface {
	AddFriend(ctx context.Context, userID, friendID string) error
	RemoveFriend(ctx context.Context, userID, friendID string) error
	ListFriends(ctx context.Context, userID string) ([]string, error)
}

type friendService struct {
	stg     storage.IFriendStorage
	log     logger.ILogger
	userStg storage.IUserStorage // agar user mavjudligini tekshirishni istasak
}

func NewFriendService(stg storage.IStorage, log logger.ILogger) FriendService {
	return &friendService{
		stg:     stg.Friend(),
		log:     log,
		userStg: stg.User(),
	}
}

func (s *friendService) AddFriend(ctx context.Context, userID, friendID string) error {
	if userID == friendID {
		return fmt.Errorf("cannot add yourself")
	}
	// 1) tekshirish: friendID user sifatida mavjudmi?
	if _, err := s.userStg.GetUserByID(ctx, friendID); err != nil {
		return fmt.Errorf("user not found")
	}
	// 2) qo'shish
	return s.stg.AddFriend(ctx, userID, friendID)
}

func (s *friendService) RemoveFriend(ctx context.Context, userID, friendID string) error {
	return s.stg.RemoveFriend(ctx, userID, friendID)
}

func (s *friendService) ListFriends(ctx context.Context, userID string) ([]string, error) {
	return s.stg.ListFriends(ctx, userID)
}

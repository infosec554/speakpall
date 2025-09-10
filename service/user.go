package service

import (
	"context"
	"errors"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/pkg/mailer"
	"speakpall/pkg/security"
	"speakpall/storage"
)

type UserService interface {
	Create(ctx context.Context, req models.SignupRequest) (string, error)
	GetForLoginByEmail(context.Context, string) (models.LoginUser, error)
	GetByID(ctx context.Context, id string) (*models.User, error)

	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	GoogleAuth(ctx context.Context, email, name, googleID string) (string, error)
}

type userService struct {
	stg        storage.IUserStorage
	log        logger.ILogger
	mailerCore *mailer.Mailer
}

func NewUserService(stg storage.IStorage, log logger.ILogger, mailerCore *mailer.Mailer) UserService {
	return &userService{
		stg:        stg.User(),
		log:        log,
		mailerCore: mailerCore,
	}
}

func (s *userService) Create(ctx context.Context, req models.SignupRequest) (string, error) {
	s.log.Info("UserService.Create called", logger.String("email", req.Email))

	id, err := s.stg.Create(ctx, req)
	if err != nil {
		s.log.Error("failed to create user", logger.Error(err))
		return "", err
	}

	s.log.Info("user successfully created", logger.String("userID", id))
	return id, nil
}

func (s *userService) GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	s.log.Info("UserService.GetForLoginByEmail called", logger.String("email", email))

	user, err := s.stg.GetForLoginByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user for login", logger.Error(err))
		return models.LoginUser{}, err
	}

	s.log.Info("user fetched for login", logger.String("userID", user.ID))
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	s.log.Info("userService.GetByID called")
	return s.stg.GetByID(ctx, id)
}
func (s *userService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	hashedOldPassword, err := s.stg.GetPasswordByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := security.CompareHashAndPassword(hashedOldPassword, oldPassword); err != nil {
		return errors.New("old password is incorrect")
	}

	hashedNewPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	return s.stg.UpdatePassword(ctx, userID, hashedNewPassword)
}

func (s *userService) GoogleAuth(ctx context.Context, email, name, googleID string) (string, error) {
	user, err := s.stg.GetForLoginByEmail(ctx, email)
	if err == nil {
		return user.ID, nil
	}
	req := models.SignupRequest{
		Name:     name,
		Email:    email,
		Password: "",
	}
	return s.stg.Create(ctx, req)
}

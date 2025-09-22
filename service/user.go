package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/pkg/mailer"
	"speakpall/pkg/security"
	"speakpall/storage"
)

type UserService interface {
	Create(ctx context.Context, req models.SignupRequest) (string, error)
	GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error)
	GetByID(ctx context.Context, id string) (*models.User, error)

	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	GoogleAuth(ctx context.Context, email, name, googleID string) (string, error)

	SetRole(ctx context.Context, userID, role string) error

	CreatePasswordResetToken(ctx context.Context, email string) (string, error)
	ValidatePasswordResetToken(ctx context.Context, token string) (string, error)
	ResetPassword(ctx context.Context, userID string, newPassword string) error
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
	s.log.Info("UserService.Create", logger.String("email", req.Email))

	// request.Password -> hash
	hashed, err := security.HashPassword(req.Password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	req.Password = hashed

	id, err := s.stg.CreateUser(ctx, req)
	if err != nil {
		s.log.Error("create user failed", logger.Error(err))
		return "", err
	}
	s.log.Info("user created", logger.String("userID", id))
	return id, nil
}

func (s *userService) GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	s.log.Info("UserService.GetForLoginByEmail", logger.String("email", email))
	u, err := s.stg.GetLoginByEmail(ctx, email)
	if err != nil {
		s.log.Error("get login by email failed", logger.Error(err))
		return models.LoginUser{}, err
	}
	return u, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	s.log.Info("UserService.GetByID", logger.String("userID", id))
	return s.stg.GetUserByID(ctx, id)
}

func (s *userService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
    user, err := s.stg.GetUserByID(ctx, userID)
    if err != nil {
        return err
    }

    if err := security.CompareHashAndPassword(user.PasswordHash, oldPassword); err != nil {
        return errors.New("old password is incorrect")
    }

    newHash, err := security.HashPassword(newPassword)
    if err != nil {
        return errors.New("failed to hash new password")
    }

    return s.stg.UpdatePasswordHash(ctx, userID, newHash)
}



func (s *userService) GoogleAuth(ctx context.Context, email, name, googleID string) (string, error) {
	// 1) bor-yo‘qligini tekshiramiz
	u, err := s.stg.GetLoginByEmail(ctx, email)
	if err == nil && u.ID != "" {
		return u.ID, nil
	}

	// 2) yo‘q bo‘lsa, random parol yaratib, hashlab create qilamiz
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", errors.New("failed to generate random password")
	}
	randomPass := base64.RawURLEncoding.EncodeToString(buf)

	hashed, err := security.HashPassword(randomPass)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	req := models.SignupRequest{
		DisplayName: name,
		Email:       email,
		Password:    hashed, // CreateUser hash kutadi (oldinda hashlangan)
	}
	id, err := s.stg.CreateUser(ctx, req)
	if err != nil {
		return "", err
	}

	// Eslatma: agar google_id ni ham saqlamoqchi bo‘lsangiz,
	// storage ga `UpdateGoogleID(userID, googleID)` kabi metod qo‘shib, shu yerda chaqirasiz.

	return id, nil
}

func (s *userService) SetRole(ctx context.Context, userID, role string) error {
	if role != models.RoleAdmin && role != models.RoleUser {
		return fmt.Errorf("invalid role: %s", role)
	}
	return s.stg.UpdateRole(ctx, userID, role)
}

func (s *userService) CreatePasswordResetToken(ctx context.Context, email string) (string, error) {
	// userni topamiz
	u, err := s.stg.GetLoginByEmail(ctx, email)
	if err != nil || u.ID == "" {
		return "", errors.New("user not found")
	}

	// 32 baytli URL-safe token generatsiya
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("failed to generate token")
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	// 1 soat amal qiladi
	expiresAt := time.Now().Add(1 * time.Hour)
	if err := s.stg.SavePasswordResetToken(ctx, u.ID, token, expiresAt); err != nil {
		return "", err
	}

	// Email yuborish (HTML)
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Password Reset</title></head>
<body style="font-family:Arial,sans-serif;background:#f4f4f4;margin:0;padding:24px">
  <div style="max-width:600px;margin:0 auto;background:#fff;padding:24px;border-radius:8px">
    <h2 style="margin:0 0 12px">Password Reset</h2>
    <p style="margin:0 0 16px">Hi, click the button below to reset your password:</p>
    <p><a href="https://yourapp.com/reset-password?token=%s"
          style="display:inline-block;background:#007bff;color:#fff;text-decoration:none;padding:12px 20px;border-radius:4px">Reset Password</a></p>
    <p style="color:#888;margin:16px 0 0">If you didn’t request this, just ignore this email.</p>
  </div>
</body></html>`, token)

	if err := s.mailerCore.Send(email, subject, body); err != nil {
		return "", errors.New("failed to send reset email")
	}

	return token, nil
}

func (s *userService) ValidatePasswordResetToken(ctx context.Context, token string) (string, error) {
	return s.stg.GetUserIDByPasswordResetToken(ctx, token, time.Now())
}

func (s *userService) ResetPassword(ctx context.Context, userID string, newPassword string) error {
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	return s.stg.UpdatePasswordHash(ctx, userID, hash)
}

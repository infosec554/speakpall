package service

import (
	// tgbotapi import qilinmoqda

	"speakpall/config"
	"speakpall/pkg/logger"
	"speakpall/pkg/mailer"
	"speakpall/storage"
)

type IServiceManager interface {
	User() UserService
	Mailer() MailerService
	Redis() RedisService

	Google() GoogleService
}

type service struct {
	userService UserService
	mailer      MailerService

	redisService  RedisService
	googleService GoogleService
}

func New(storage storage.IStorage, log logger.ILogger, mailerCore *mailer.Mailer, redis storage.IRedisStorage, googleCfg config.OAuthProviderConfig) IServiceManager {
	return &service{
		userService: NewUserService(storage, log, mailerCore),
		mailer:      NewMailerService(mailerCore),

		redisService:  NewRedisService(redis, log),
		googleService: NewGoogleService(GoogleOAuthConfig(googleCfg)), // <-- config ni uzatish!

	}
}

func (s *service) User() UserService {
	return s.userService
}

func (s *service) Mailer() MailerService {
	return s.mailer
}

func (s *service) Redis() RedisService {
	return s.redisService
}

func (s *service) Google() GoogleService {
	return s.googleService
}

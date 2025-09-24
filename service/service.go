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
	Profile() ProfileService
	Settings() SettingsService
	Matchs() MatchsService
	Interes() InteresService
	Friend() FriendService
}

type service struct {
	userService UserService
	mailer      MailerService

	redisService    RedisService
	googleService   GoogleService
	profileService  ProfileService
	settingsService SettingsService
	matchsService   MatchsService
	interesService  InteresService
	friendService  FriendService
}

func New(storage storage.IStorage, log logger.ILogger, mailerCore *mailer.Mailer, redis storage.IRedisStorage, googleCfg config.OAuthProviderConfig) IServiceManager {
	return &service{
		userService: NewUserService(storage, log, mailerCore),
		mailer:      NewMailerService(mailerCore),

		redisService:    NewRedisService(redis, log),
		googleService:   NewGoogleService(GoogleOAuthConfig(googleCfg)), // <-- config ni uzatish!
		profileService:  NewProfileService(storage, log),
		settingsService: NewSettingsService(storage, log),
		matchsService:   NewMatchsService(storage, log),
		interesService: NewInteresService(storage,log),
		friendService: NewFriendService(storage,log),
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

func (s *service) Profile() ProfileService {
	return s.profileService
}
func (s *service) Settings() SettingsService {
	return s.settingsService
}
func (s *service) Matchs() MatchsService {
	return s.matchsService
}

func (s *service) Interes() InteresService  {
	return s.interesService
}

func (s *service) Friend()  FriendService {
	return s.friendService
}
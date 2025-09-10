package main

import (
	"context"

	"speakpall/api"
	"speakpall/config"
	"speakpall/pkg/logger"
	"speakpall/pkg/mailer"
	"speakpall/service"
	"speakpall/storage/postgres"
	"speakpall/storage/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	pgStore, err := postgres.New(context.Background(), cfg, log, nil)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return
	}
	defer pgStore.Close()

	mailService := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPSenderName)
	redisStore := redis.New(cfg)

	services := service.New(pgStore, log, mailService, redisStore, cfg.Google)

	server := api.New(services, log)
	log.Info("Service is running on", logger.Int("port", 8081))
	if err = server.Run("localhost:8011"); err != nil {
		panic(err)
	}
}

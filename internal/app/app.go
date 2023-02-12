package app

import (
	"context"
	"github.com/serhiimazurok/ecoflow-status-bot/internal/bot"
	"github.com/serhiimazurok/ecoflow-status-bot/internal/config"
	"github.com/serhiimazurok/ecoflow-status-bot/internal/repository"
	"github.com/serhiimazurok/ecoflow-status-bot/internal/service"
	"github.com/serhiimazurok/ecoflow-status-bot/pkg/database/mongodb"
	"github.com/serhiimazurok/ecoflow-status-bot/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logger.Error(err)

		return
	}

	// Dependencies
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)

		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	//services
	_ = service.NewServices(service.Deps{
		Repos: repos,
	})

	// Telegram Bot
	tg, err := bot.New(cfg)
	if err != nil {
		logger.Error(err)

		return
	}

	go tg.Run()

	logger.Info("Bot started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	tg.Stop()

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}

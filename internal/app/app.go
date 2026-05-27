package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"awesomeProject/internal/config"
	"awesomeProject/internal/database"
	"awesomeProject/internal/handler"
	httpserver "awesomeProject/internal/http"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
)

type App struct {
	config config.Config
}

func New(config config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() error {
	ctx := context.Background()

	db, err := database.Connect(ctx, a.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("get underlying db err:", err)
			return
		}

		if err := sqlDB.Close(); err != nil {
			log.Println("close db err:", err)
		}
	}()

	slog.Info("connect to database")

	webhookEventRepository := repository.NewPostgresWebhookEventRepository(db)
	webhookEventService := service.NewWebhookEventService(webhookEventRepository)
	webhookEventHandler := handler.NewWebhookEventHandler(webhookEventService)

	router := httpserver.NewRouter(httpserver.RouterDependencies{
		WebhookEventHandler: webhookEventHandler,
	})

	return router.Run(":" + a.config.Port)
}

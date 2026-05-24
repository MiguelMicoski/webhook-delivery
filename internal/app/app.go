package app

import (
	"awesomeProject/internal/config"
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
	webhookEventRepository := repository.NewMemoryWebhookEventRepository()
	webhookEventService := service.NewWebhookEventService(webhookEventRepository)
	webhookEventHandler := handler.NewWebhookEventHandler(webhookEventService)

	router := httpserver.NewRouter(httpserver.RouterDependencies{
		WebhookEventHandler: webhookEventHandler,
	})

	return router.Run(":" + a.config.Port)
}

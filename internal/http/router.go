package http

import (
	"awesomeProject/internal/handler"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	WebhookEventHandler *handler.WebhookEventHandler
}

func NewRouter(dependencies RouterDependencies) *gin.Engine {
	router := gin.Default()

	router.POST("/webhook-events", dependencies.WebhookEventHandler.Create)
	router.GET("/webhook-events/:id", dependencies.WebhookEventHandler.FindByID)

	return router
}

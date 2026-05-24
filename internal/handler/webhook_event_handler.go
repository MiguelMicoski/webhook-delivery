package handler

import (
	"errors"
	"net/http"

	"awesomeProject/internal/service"

	"github.com/gin-gonic/gin"
)

type WebhookEventHandler struct {
	webhookEventService *service.WebhookEventService
}

type createWebhookEventRequest struct {
	TargetURL string         `json:"target_url" binding:"required"`
	Payload   map[string]any `json:"payload" binding:"required"`
}

func NewWebhookEventHandler(webhookEventService *service.WebhookEventService) *WebhookEventHandler {
	return &WebhookEventHandler{
		webhookEventService: webhookEventService,
	}
}

func (h *WebhookEventHandler) Create(ctx *gin.Context) {
	var request createWebhookEventRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	event, err := h.webhookEventService.Create(ctx.Request.Context(), service.CreateWebhookEventInput{
		TargetURL: request.TargetURL,
		Payload:   request.Payload,
	})
	if err != nil {
		status := http.StatusInternalServerError
		message := "internal server error"

		if errors.Is(err, service.ErrTargetURLRequired) ||
			errors.Is(err, service.ErrInvalidTargetURL) ||
			errors.Is(err, service.ErrPayloadRequired) {
			status = http.StatusBadRequest
			message = err.Error()
		}

		ctx.JSON(status, gin.H{
			"error": message,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": event,
	})
}

func (h *WebhookEventHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	event, err := h.webhookEventService.FindByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrWebhookEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": event,
	})
}

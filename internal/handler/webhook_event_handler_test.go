package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"

	"github.com/gin-gonic/gin"
)

func TestWebhookEventHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookEventRepository := repository.NewMemoryWebhookEventRepository()
	webhookEventService := service.NewWebhookEventService(webhookEventRepository)
	webhookEventHandler := NewWebhookEventHandler(webhookEventService)

	router := gin.New()
	router.POST("/webhook-events", webhookEventHandler.Create)

	body := `{"target_url":"https://example.com/webhook","payload":{"type":"order.created"}}`
	request := httptest.NewRequest(http.MethodPost, "/webhook-events", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, response.Code, response.Body.String())
	}

	if !strings.Contains(response.Body.String(), `"status":"pending"`) {
		t.Fatalf("expected response body to contain pending status, got %s", response.Body.String())
	}
}

func TestWebhookEventHandlerFindByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookEventRepository := repository.NewMemoryWebhookEventRepository()
	webhookEventService := service.NewWebhookEventService(webhookEventRepository)
	webhookEventHandler := NewWebhookEventHandler(webhookEventService)

	createdEvent, err := webhookEventService.Create(t.Context(), service.CreateWebhookEventInput{
		TargetURL: "https://example.com/webhook",
		Payload: map[string]any{
			"type": "order.created",
		},
	})
	if err != nil {
		t.Fatalf("expected no error creating event, got %v", err)
	}

	router := gin.New()
	router.GET("/webhook-events/:id", webhookEventHandler.FindByID)

	request := httptest.NewRequest(http.MethodGet, "/webhook-events/"+createdEvent.ID, nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, response.Code, response.Body.String())
	}

	if !strings.Contains(response.Body.String(), `"id":"`+createdEvent.ID+`"`) {
		t.Fatalf("expected response body to contain created event ID, got %s", response.Body.String())
	}
}

func TestWebhookEventHandlerFindByIDReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookEventRepository := repository.NewMemoryWebhookEventRepository()
	webhookEventService := service.NewWebhookEventService(webhookEventRepository)
	webhookEventHandler := NewWebhookEventHandler(webhookEventService)

	router := gin.New()
	router.GET("/webhook-events/:id", webhookEventHandler.FindByID)

	request := httptest.NewRequest(http.MethodGet, "/webhook-events/evt_missing", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusNotFound, response.Code, response.Body.String())
	}
}

func TestWebhookEventHandlerCreateRejectsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webhookEventRepository := repository.NewMemoryWebhookEventRepository()
	webhookEventService := service.NewWebhookEventService(webhookEventRepository)
	webhookEventHandler := NewWebhookEventHandler(webhookEventService)

	router := gin.New()
	router.POST("/webhook-events", webhookEventHandler.Create)

	request := httptest.NewRequest(http.MethodPost, "/webhook-events", strings.NewReader(`{"target_url":`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

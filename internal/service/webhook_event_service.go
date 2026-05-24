package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"time"

	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
)

var (
	ErrTargetURLRequired    = errors.New("target_url is required")
	ErrInvalidTargetURL     = errors.New("target_url must be a valid http or https URL")
	ErrPayloadRequired      = errors.New("payload is required")
	ErrWebhookEventNotFound = errors.New("webhook event not found")
)

type CreateWebhookEventInput struct {
	TargetURL string
	Payload   map[string]any
}

type WebhookEventService struct {
	repository repository.WebhookEventRepository
}

func NewWebhookEventService(repository repository.WebhookEventRepository) *WebhookEventService {
	return &WebhookEventService{
		repository: repository,
	}
}

func (s *WebhookEventService) Create(ctx context.Context, input CreateWebhookEventInput) (model.WebhookEvent, error) {
	if input.TargetURL == "" {
		return model.WebhookEvent{}, ErrTargetURLRequired
	}

	if !isHTTPURL(input.TargetURL) {
		return model.WebhookEvent{}, ErrInvalidTargetURL
	}

	if input.Payload == nil {
		return model.WebhookEvent{}, ErrPayloadRequired
	}

	event := model.WebhookEvent{
		ID:        newWebhookEventID(),
		TargetURL: input.TargetURL,
		Payload:   input.Payload,
		Status:    model.WebhookEventStatusPending,
		CreatedAt: time.Now().UTC(),
	}

	createdEvent, err := s.repository.Create(ctx, event)
	if err != nil {
		return model.WebhookEvent{}, fmt.Errorf("create webhook event: %w", err)
	}

	return createdEvent, nil
}

func (s *WebhookEventService) FindByID(ctx context.Context, id string) (model.WebhookEvent, error) {
	event, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWebhookEventNotFound) {
			return model.WebhookEvent{}, ErrWebhookEventNotFound
		}

		return model.WebhookEvent{}, fmt.Errorf("find webhook event by id: %w", err)
	}

	return event, nil
}

func isHTTPURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func newWebhookEventID() string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return fmt.Sprintf("evt_%d", time.Now().UnixNano())
	}

	return "evt_" + hex.EncodeToString(bytes[:])
}

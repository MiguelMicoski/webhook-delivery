package service

import (
	"context"
	"errors"
	"testing"

	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
)

type fakeWebhookEventRepository struct {
	event model.WebhookEvent
	err   error
}

func (r *fakeWebhookEventRepository) Create(ctx context.Context, event model.WebhookEvent) (model.WebhookEvent, error) {
	if r.err != nil {
		return model.WebhookEvent{}, r.err
	}

	r.event = event

	return event, nil
}

func (r *fakeWebhookEventRepository) FindByID(ctx context.Context, id string) (model.WebhookEvent, error) {
	if r.err != nil {
		return model.WebhookEvent{}, r.err
	}

	if r.event.ID != id {
		return model.WebhookEvent{}, repository.ErrWebhookEventNotFound
	}

	return r.event, nil
}

type fakeWebhookEventPublisher struct {
	eventID string
	err     error
}

func (p *fakeWebhookEventPublisher) PublishWebhookEventCreated(ctx context.Context, eventID string) error {
	if p.err != nil {
		return p.err
	}

	p.eventID = eventID

	return nil
}

func TestWebhookEventServiceCreate(t *testing.T) {
	repository := &fakeWebhookEventRepository{}
	publisher := &fakeWebhookEventPublisher{}
	service := NewWebhookEventService(repository, publisher)

	event, err := service.Create(context.Background(), CreateWebhookEventInput{
		TargetURL: "https://example.com/webhook",
		Payload: map[string]any{
			"type": "order.created",
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if event.ID == "" {
		t.Fatal("expected event ID to be filled")
	}

	if event.Status != model.WebhookEventStatusPending {
		t.Fatalf("expected status %q, got %q", model.WebhookEventStatusPending, event.Status)
	}

	if repository.event.ID != event.ID {
		t.Fatal("expected repository to receive created event")
	}

	if publisher.eventID != event.ID {
		t.Fatal("expected publisher to receive created event ID")
	}
}

func TestWebhookEventServiceCreateValidatesTargetURL(t *testing.T) {
	service := NewWebhookEventService(&fakeWebhookEventRepository{}, nil)

	_, err := service.Create(context.Background(), CreateWebhookEventInput{
		TargetURL: "not-a-url",
		Payload: map[string]any{
			"type": "order.created",
		},
	})
	if !errors.Is(err, ErrInvalidTargetURL) {
		t.Fatalf("expected ErrInvalidTargetURL, got %v", err)
	}
}

func TestWebhookEventServiceFindByID(t *testing.T) {
	repository := &fakeWebhookEventRepository{
		event: model.WebhookEvent{
			ID:        "evt_123",
			TargetURL: "https://example.com/webhook",
			Payload: map[string]any{
				"type": "order.created",
			},
			Status: model.WebhookEventStatusPending,
		},
	}
	service := NewWebhookEventService(repository, nil)

	event, err := service.FindByID(context.Background(), "evt_123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if event.ID != "evt_123" {
		t.Fatalf("expected event ID %q, got %q", "evt_123", event.ID)
	}
}

func TestWebhookEventServiceFindByIDReturnsNotFound(t *testing.T) {
	service := NewWebhookEventService(&fakeWebhookEventRepository{}, nil)

	_, err := service.FindByID(context.Background(), "evt_missing")
	if !errors.Is(err, ErrWebhookEventNotFound) {
		t.Fatalf("expected ErrWebhookEventNotFound, got %v", err)
	}
}

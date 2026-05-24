package repository

import (
	"context"
	"errors"
	"sync"

	"awesomeProject/internal/model"
)

var ErrWebhookEventNotFound = errors.New("webhook event not found")

type WebhookEventRepository interface {
	Create(ctx context.Context, event model.WebhookEvent) (model.WebhookEvent, error)
	FindByID(ctx context.Context, id string) (model.WebhookEvent, error)
}

type MemoryWebhookEventRepository struct {
	mu     sync.RWMutex
	events map[string]model.WebhookEvent
}

func NewMemoryWebhookEventRepository() *MemoryWebhookEventRepository {
	return &MemoryWebhookEventRepository{
		events: make(map[string]model.WebhookEvent),
	}
}

func (r *MemoryWebhookEventRepository) Create(ctx context.Context, event model.WebhookEvent) (model.WebhookEvent, error) {
	if err := ctx.Err(); err != nil {
		return model.WebhookEvent{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[event.ID] = event

	return event, nil
}

func (r *MemoryWebhookEventRepository) FindByID(ctx context.Context, id string) (model.WebhookEvent, error) {
	if err := ctx.Err(); err != nil {
		return model.WebhookEvent{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	event, ok := r.events[id]
	if !ok {
		return model.WebhookEvent{}, ErrWebhookEventNotFound
	}

	return event, nil
}

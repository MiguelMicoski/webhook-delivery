package repository

import (
	"context"
	"errors"
	"fmt"

	"awesomeProject/internal/model"

	"gorm.io/gorm"
)

type PostgresWebhookEventRepository struct {
	db *gorm.DB
}

func NewPostgresWebhookEventRepository(db *gorm.DB) *PostgresWebhookEventRepository {
	return &PostgresWebhookEventRepository{
		db: db,
	}
}

func (r *PostgresWebhookEventRepository) Create(ctx context.Context, event model.WebhookEvent) (model.WebhookEvent, error) {
	if err := r.db.WithContext(ctx).Create(&event).Error; err != nil {
		return model.WebhookEvent{}, fmt.Errorf("insert webhook event: %w", err)
	}

	return event, nil
}

func (r *PostgresWebhookEventRepository) FindByID(ctx context.Context, id string) (model.WebhookEvent, error) {
	var event model.WebhookEvent

	err := r.db.WithContext(ctx).First(&event, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.WebhookEvent{}, ErrWebhookEventNotFound
		}
		return model.WebhookEvent{}, fmt.Errorf("select webhook event by id: %w", err)
	}

	return event, nil
}

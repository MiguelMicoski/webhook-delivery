package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"awesomeProject/internal/model"
)

type PostgresWebhookEventRepository struct {
	db *sql.DB
}

func NewPostgresWebhookEventRepository(db *sql.DB) *PostgresWebhookEventRepository {
	return &PostgresWebhookEventRepository{
		db: db,
	}
}

func (r *PostgresWebhookEventRepository) Create(ctx context.Context, event model.WebhookEvent) (model.WebhookEvent, error) {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return model.WebhookEvent{}, fmt.Errorf("marshal webhook event payload: %w", err)
	}

	const query = `
		INSERT INTO webhook_events (id, target_url, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err := r.db.ExecContext(ctx, query, event.ID, event.TargetURL, payload, event.Status, event.CreatedAt); err != nil {
		return model.WebhookEvent{}, fmt.Errorf("insert webhook event: %w", err)
	}

	return event, nil
}

func (r *PostgresWebhookEventRepository) FindByID(ctx context.Context, id string) (model.WebhookEvent, error) {
	const query = `
		SELECT id, target_url, payload, status, created_at
		FROM webhook_events
		WHERE id = $1
	`

	var event model.WebhookEvent
	var payload []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.TargetURL,
		&payload,
		&event.Status,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.WebhookEvent{}, ErrWebhookEventNotFound
		}

		return model.WebhookEvent{}, fmt.Errorf("select webhook event by id: %w", err)
	}

	if err := json.Unmarshal(payload, &event.Payload); err != nil {
		return model.WebhookEvent{}, fmt.Errorf("unmarshal webhook event payload: %w", err)
	}

	return event, nil
}

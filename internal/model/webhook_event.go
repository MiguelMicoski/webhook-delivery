package model

import "time"

type WebhookEventStatus string

const (
	WebhookEventStatusPending WebhookEventStatus = "pending"
)

type WebhookEvent struct {
	ID        string             `json:"id"`
	TargetURL string             `json:"target_url"`
	Payload   map[string]any     `json:"payload"`
	Status    WebhookEventStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
}

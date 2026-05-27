package model

import "time"

type WebhookEventStatus string

const (
	WebhookEventStatusPending WebhookEventStatus = "pending"
)

type WebhookEvent struct {
	ID        string             `json:"id" gorm:"primaryKey;column:id"`
	TargetURL string             `json:"target_url" gorm:"column:target_url;not null"`
	Payload   map[string]any     `json:"payload" gorm:"column:payload;type:jsonb;serializer:json;not null"`
	Status    WebhookEventStatus `json:"status" gorm:"column:status;not null"`
	CreatedAt time.Time          `json:"created_at" gorm:"column:created_at;not null"`
}

func (WebhookEvent) TableName() string {
	return "webhook_events"
}

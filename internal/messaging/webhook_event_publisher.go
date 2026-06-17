package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WebhookEventPublisher struct {
	channel    *amqp.Channel
	confirms   <-chan amqp.Confirmation
	exchange   string
	routingKey string
	mu         sync.Mutex
}

type webhookEventCreatedMessage struct {
	EventID string `json:"event_id"`
}

func NewWebhookEventPublisher(rabbitmq *RabbitMQ) *WebhookEventPublisher {
	config := rabbitmq.Config()

	return &WebhookEventPublisher{
		channel:    rabbitmq.Channel(),
		confirms:   rabbitmq.Channel().NotifyPublish(make(chan amqp.Confirmation, 1)),
		exchange:   config.Exchange,
		routingKey: config.RoutingKey,
	}
}

func (p *WebhookEventPublisher) PublishWebhookEventCreated(ctx context.Context, eventID string) error {
	body, err := json.Marshal(webhookEventCreatedMessage{
		EventID: eventID,
	})
	if err != nil {
		return fmt.Errorf("marshal webhook event created message: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.channel.PublishWithContext(
		ctx,
		p.exchange,
		p.routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("publish webhook event created message: %w", err)
	}

	select {
	case confirmation := <-p.confirms:
		if !confirmation.Ack {
			return fmt.Errorf("rabbitmq did not confirm webhook event created message")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

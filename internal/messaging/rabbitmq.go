package messaging

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL                string
	Exchange           string
	Queue              string
	RoutingKey         string
	DeadLetterExchange string
	DeadLetterQueue    string
}

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	config     RabbitMQConfig
}

func ConnectRabbitMQ(ctx context.Context, config RabbitMQConfig) (*RabbitMQ, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}

	connection, err := amqp.DialConfig(config.URL, amqp.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		_ = connection.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	if err := channel.Confirm(false); err != nil {
		_ = channel.Close()
		_ = connection.Close()
		return nil, fmt.Errorf("enable publisher confirms: %w", err)
	}

	rabbitmq := &RabbitMQ{
		connection: connection,
		channel:    channel,
		config:     config,
	}

	if err := rabbitmq.declareTopology(ctx); err != nil {
		_ = rabbitmq.Close()
		return nil, err
	}

	return rabbitmq, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		_ = r.connection.Close()
		return fmt.Errorf("close rabbitmq channel: %w", err)
	}

	if err := r.connection.Close(); err != nil {
		return fmt.Errorf("close rabbitmq connection: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQ) Config() RabbitMQConfig {
	return r.config
}

func (r *RabbitMQ) declareTopology(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := r.channel.ExchangeDeclare(
		r.config.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	if err := r.channel.ExchangeDeclare(
		r.config.DeadLetterExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare dead-letter exchange: %w", err)
	}

	if _, err := r.channel.QueueDeclare(
		r.config.DeadLetterQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare dead-letter queue: %w", err)
	}

	if err := r.channel.QueueBind(
		r.config.DeadLetterQueue,
		r.config.RoutingKey,
		r.config.DeadLetterExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind dead-letter queue: %w", err)
	}

	if _, err := r.channel.QueueDeclare(
		r.config.Queue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    r.config.DeadLetterExchange,
			"x-dead-letter-routing-key": r.config.RoutingKey,
		},
	); err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	if err := r.channel.QueueBind(
		r.config.Queue,
		r.config.RoutingKey,
		r.config.Exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	return nil
}

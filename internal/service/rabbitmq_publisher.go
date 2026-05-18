package service

import (
	"context"
	"encoding/json"
	"fmt"

	"GoFrioCalor/config"
	"GoFrioCalor/internal/dto"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type RabbitMQPublisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config *config.RabbitMQConfig
}

func NewRabbitMQPublisher(rabbitConfig *config.RabbitMQConfig) (*RabbitMQPublisher, error) {
	conn, err := config.ConnectRabbitMQ(rabbitConfig)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = config.DeclareQueue(ch, rabbitConfig.Queue)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Info().
		Str("queue", rabbitConfig.Queue).
		Str("host", rabbitConfig.Host).
		Msg("RabbitMQ Publisher connected successfully")

	return &RabbitMQPublisher{
		conn:   conn,
		ch:     ch,
		config: rabbitConfig,
	}, nil
}

func (p *RabbitMQPublisher) PublishWorkOrder(ctx context.Context, message dto.WorkOrderMessageDTO) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.ch.PublishWithContext(
		ctx,
		"",
		p.config.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Error().
			Err(err).
			Int("delivery_id", message.DeliveryID).
			Str("nro_cta", message.NroCta).
			Msg("Failed to publish work order message")
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Info().
		Int("delivery_id", message.DeliveryID).
		Str("nro_cta", message.NroCta).
		Str("queue", p.config.Queue).
		Msg("Work order message published successfully")

	return nil
}

func (p *RabbitMQPublisher) Close() error {
	if err := p.ch.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing RabbitMQ channel")
	}
	if err := p.conn.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing RabbitMQ connection")
		return err
	}
	log.Info().Msg("RabbitMQ Publisher closed")
	return nil
}

func (p *RabbitMQPublisher) IsHealthy() bool {
	return p.conn != nil && !p.conn.IsClosed()
}

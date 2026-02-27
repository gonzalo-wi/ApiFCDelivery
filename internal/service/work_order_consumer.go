package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"GoFrioCalor/config"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// WorkOrderConsumer procesa mensajes de RabbitMQ para crear órdenes de trabajo
type WorkOrderConsumer struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	config         *config.RabbitMQConfig
	workOrderStore store.WorkOrderStore
	deliveryStore  store.DeliveryStore
	pdfGenerator   WorkOrderPDFGenerator
	emailService   EmailService
	stopChan       chan struct{}
}

// NewWorkOrderConsumer crea un nuevo consumidor
func NewWorkOrderConsumer(
	rabbitConfig *config.RabbitMQConfig,
	workOrderStore store.WorkOrderStore,
	deliveryStore store.DeliveryStore,
	pdfGenerator WorkOrderPDFGenerator,
	emailService EmailService,
) (*WorkOrderConsumer, error) {
	conn, err := config.ConnectRabbitMQ(rabbitConfig)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declarar la cola
	_, err = config.DeclareQueue(ch, rabbitConfig.Queue)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Configurar QoS - procesar un mensaje a la vez
	err = ch.Qos(1, 0, false)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	log.Info().
		Str("queue", rabbitConfig.Queue).
		Str("host", rabbitConfig.Host).
		Msg("Work Order Consumer connected successfully")

	return &WorkOrderConsumer{
		conn:           conn,
		ch:             ch,
		config:         rabbitConfig,
		workOrderStore: workOrderStore,
		deliveryStore:  deliveryStore,
		pdfGenerator:   pdfGenerator,
		emailService:   emailService,
		stopChan:       make(chan struct{}),
	}, nil
}

// Start inicia el consumo de mensajes
func (c *WorkOrderConsumer) Start(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		c.config.Queue, // queue
		"",             // consumer
		false,          // auto-ack (false para confirmar manualmente)
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Info().Msg("Work Order Consumer started. Waiting for messages...")

	go func() {
		for {
			select {
			case <-c.stopChan:
				log.Info().Msg("Work Order Consumer stopped")
				return
			case d, ok := <-msgs:
				if !ok {
					log.Warn().Msg("Channel closed, stopping consumer")
					return
				}
				c.processMessage(ctx, d)
			}
		}
	}()

	return nil
}

// processMessage procesa un mensaje individual
func (c *WorkOrderConsumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	log.Info().
		Str("message_id", msg.MessageId).
		Int("body_size", len(msg.Body)).
		Msg("Processing work order message")

	var workOrderMsg dto.WorkOrderMessageDTO
	err := json.Unmarshal(msg.Body, &workOrderMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal message")
		// Rechazar mensaje - no será reprocesado
		msg.Nack(false, false)
		return
	}

	// Procesar la orden de trabajo
	err = c.createWorkOrder(ctx, workOrderMsg)
	if err != nil {
		log.Error().
			Err(err).
			Int("delivery_id", workOrderMsg.DeliveryID).
			Msg("Failed to process work order")

		// Rechazar el mensaje y devolverlo a la cola para reintento
		// TODO: Implementar dead letter queue después de N reintentos
		msg.Nack(false, true)
		return
	}

	// Confirmar el mensaje
	msg.Ack(false)
	log.Info().
		Int("delivery_id", workOrderMsg.DeliveryID).
		Msg("Work order processed successfully")
}

// createWorkOrder crea la orden de trabajo, genera PDF y envía email
func (c *WorkOrderConsumer) createWorkOrder(ctx context.Context, msg dto.WorkOrderMessageDTO) error {
	// 1. Generar número de orden
	orderNumber, err := c.workOrderStore.GetNextOrderNumber(ctx)
	if err != nil {
		return fmt.Errorf("error generando número de orden: %w", err)
	}

	// 2. Crear orden de trabajo
	workOrder := &models.WorkOrder{
		OrderNumber: orderNumber,
		NroCta:      msg.NroCta,
		NroRto:      msg.NroRto,
		Name:        msg.Name,
		Address:     msg.Address,
		Localidad:   msg.Locality,
		TipoAccion:  msg.TipoAccion,
		CreatedAt:   time.Now(),
	}

	err = c.workOrderStore.Create(ctx, workOrder)
	if err != nil {
		return fmt.Errorf("error creando orden de trabajo: %w", err)
	}

	log.Info().
		Str("order_number", orderNumber).
		Int("work_order_id", workOrder.ID).
		Msg("Work order created")

	// 3. Generar PDF (si el servicio está disponible)
	var pdfPath string
	if c.pdfGenerator != nil {
		pdfPath, err = c.pdfGenerator.GenerateWorkOrderPDF(ctx, workOrder, msg.Dispensers)
		if err != nil {
			log.Error().Err(err).Msg("Error generando PDF, continuando sin él")
			// No fallar si el PDF falla
		} else {
			log.Info().Str("pdf_path", pdfPath).Msg("PDF generated successfully")
		}
	}

	// 4. Enviar email (si el servicio está disponible)
	if c.emailService != nil {
		err = c.emailService.SendWorkOrderEmail(ctx, workOrder, pdfPath)
		if err != nil {
			log.Error().Err(err).Msg("Error enviando email, continuando")
			// No fallar si el email falla
		} else {
			log.Info().Str("email", "cliente@example.com").Msg("Email sent successfully")
		}
	}

	// 5. TODO: Guardar PDF en storage (implementar cuando esté disponible)
	// c.storageService.Upload(pdfPath)

	log.Info().
		Int("delivery_id", msg.DeliveryID).
		Str("order_number", orderNumber).
		Msg("Work order workflow completed")

	return nil
}

// Stop detiene el consumidor
func (c *WorkOrderConsumer) Stop() error {
	close(c.stopChan)

	if err := c.ch.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing RabbitMQ channel")
	}
	if err := c.conn.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing RabbitMQ connection")
		return err
	}
	log.Info().Msg("Work Order Consumer closed")
	return nil
}

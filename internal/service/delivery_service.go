package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/rs/zerolog/log"
)

type DeliveryService interface {
	FindAll(ctx context.Context) ([]models.Delivery, error)
	FindByID(ctx context.Context, id int) (*models.Delivery, error)
	FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error)
	FindByRto(ctx context.Context, nroRto string, fechaAccion *time.Time) ([]models.Delivery, error)
	FindByFechaAccion(ctx context.Context, fecha string) ([]models.Delivery, error)
	FindByFechaAndNroCta(ctx context.Context, fechaAccion, nroCta string) (*models.Delivery, error)
	Create(ctx context.Context, delivery *models.Delivery) error
	Update(ctx context.Context, delivery *models.Delivery) error
	Delete(ctx context.Context, id int) error
	CreateFromInfobip(ctx context.Context, req dto.InfobipDeliveryRequest) (*models.Delivery, error)
}
type deliveryService struct {
	store        store.DeliveryStore
	emailService EmailService
}

func NewDeliveryService(store store.DeliveryStore) DeliveryService {
	return &deliveryService{
		store:        store,
		emailService: nil,
	}
}

func NewDeliveryServiceWithEmail(store store.DeliveryStore, emailService EmailService) DeliveryService {
	return &deliveryService{
		store:        store,
		emailService: emailService,
	}
}

func (s *deliveryService) FindAll(ctx context.Context) ([]models.Delivery, error) {
	return s.store.FindAll(ctx)
}

func (s *deliveryService) FindByID(ctx context.Context, id int) (*models.Delivery, error) {
	return s.store.FindByID(ctx, id)
}

func (s *deliveryService) FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error) {
	return s.store.FindByFilters(ctx, nroCta, fechaAccion)
}

func (s *deliveryService) FindByRto(ctx context.Context, rto string, fechaAccion *time.Time) ([]models.Delivery, error) {
	return s.store.FindByRto(ctx, rto, fechaAccion)
}

func (s *deliveryService) FindByFechaAccion(ctx context.Context, fecha string) ([]models.Delivery, error) {
	return s.store.FindByFechaAccion(ctx, fecha)
}

func (s *deliveryService) FindByFechaAndNroCta(ctx context.Context, fechaAccion, nroCta string) (*models.Delivery, error) {
	return s.store.FindByFechaAndNroCta(ctx, fechaAccion, nroCta)
}

// Al momento de crear la entrega se genera el token para el cliente
func (s *deliveryService) Create(ctx context.Context, delivery *models.Delivery) error {
	delivery.Token = s.generateToken()
	if delivery.FechaAccion.IsZero() {
		delivery.FechaAccion = models.CustomDate{Time: time.Now()}
	}
	return s.store.Create(ctx, delivery)
}

func (s *deliveryService) Update(ctx context.Context, delivery *models.Delivery) error {
	return s.store.Update(ctx, delivery)
}

func (s *deliveryService) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}

// CreateFromInfobip crea una entrega desde el chatbot de Infobip
// Implementa idempotencia: si ya existe una entrega con el mismo session_id, la devuelve
// Maneja concurrencia de forma segura mediante índice único en BD y transacciones
func (s *deliveryService) CreateFromInfobip(ctx context.Context, req dto.InfobipDeliveryRequest) (*models.Delivery, error) {
	// Verificar si ya existe una entrega con este session_id (idempotencia)
	if req.SessionID != "" {
		existingDelivery, err := s.store.FindBySessionID(ctx, req.SessionID)
		if err != nil {
			return nil, fmt.Errorf("error verificando session_id existente: %w", err)
		}
		if existingDelivery != nil {
			// Ya existe una entrega con este session_id, devolverla (idempotente)
			return existingDelivery, nil
		}
	}

	// Validar cantidad total de dispensers
	cantidadTotal := req.Tipos.P + req.Tipos.M
	if err := validateDispenserQuantity(cantidadTotal); err != nil {
		return nil, err
	}

	// Parsear fecha de acción
	fechaAccion, err := parseFechaAccion(req.FechaAccion)
	if err != nil {
		return nil, err
	}

	// Crear items de dispensers
	itemDispensers := createItemDispensers(req.Tipos.P, req.Tipos.M)

	// Preparar SessionID como puntero (nil si está vacío)
	var sessionIDPtr *string
	if req.SessionID != "" {
		sessionIDPtr = &req.SessionID
	}

	// Crear la entrega
	delivery := &models.Delivery{
		NroCta:         req.NroCta,
		NroRto:         req.NroRto,
		Email:          req.Email,
		ItemDispensers: itemDispensers,
		Cantidad:       cantidadTotal,
		Estado:         models.Pendiente,
		TipoEntrega:    req.TipoEntrega,
		EntregadoPor:   req.EntregadoPor,
		SessionID:      sessionIDPtr,
		FechaAccion:    fechaAccion,
	}

	// Generar token de 4 dígitos (thread-safe)
	delivery.Token = s.generateToken()

	// Guardar en base de datos (la concurrencia se maneja a nivel de BD)
	if err := s.store.Create(ctx, delivery); err != nil {
		return nil, fmt.Errorf("error creando entrega: %w", err)
	}

	// NO enviar email aquí, se enviará cuando se complete la entrega

	return delivery, nil
}

// generar token de 4 digitos
func (s *deliveryService) generateToken() string {
	rangeSize := int64(constants.TOKEN_MAX - constants.TOKEN_MIN + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
	}
	token := int(n.Int64()) + constants.TOKEN_MIN
	return fmt.Sprintf("%04d", token)
}

// sendDeliveryConfirmationEmail envía un email de confirmación de entrega
func (s *deliveryService) sendDeliveryConfirmationEmail(ctx context.Context, delivery *models.Delivery) {
	subject := fmt.Sprintf("Confirmación de Entrega - Token: %s", delivery.Token)

	htmlBody := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px;">
					✅ Confirmación de Entrega Programada
				</h2>
				
				<p>Estimado cliente,</p>
				
				<p>Su entrega ha sido programada exitosamente. A continuación los detalles:</p>
				
				<div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
					<p style="margin: 5px 0;"><strong>🔑 Token de Validación:</strong> <span style="font-size: 24px; color: #e74c3c; font-weight: bold;">%s</span></p>
					<p style="margin: 5px 0;"><strong>📦 Cuenta:</strong> %s</p>
					<p style="margin: 5px 0;"><strong>🚚 Ruta:</strong> %s</p>
					<p style="margin: 5px 0;"><strong>📅 Fecha:</strong> %s</p>
					<p style="margin: 5px 0;"><strong>📊 Cantidad de Dispensers:</strong> %d</p>
					<p style="margin: 5px 0;"><strong>🔧 Tipo de Entrega:</strong> %s</p>
				</div>
				
				<div style="background-color: #fff3cd; padding: 15px; border-left: 4px solid #ffc107; margin: 20px 0;">
					<p style="margin: 0;"><strong>⚠️ Importante:</strong> Tenga este token a mano cuando llegue el repartidor. Será necesario para validar y completar la entrega.</p>
				</div>
				
				<p style="color: #7f8c8d; font-size: 12px; margin-top: 30px; border-top: 1px solid #ecf0f1; padding-top: 15px;">
					Este es un email automático. Por favor no responda a este mensaje.<br>
					<strong>El Jumillano - Sistema de Gestión de Entregas</strong>
				</p>
			</div>
		</body>
		</html>
	`,
		delivery.Token,
		delivery.NroCta,
		delivery.NroRto,
		delivery.FechaAccion.Format("02/01/2006"),
		delivery.Cantidad,
		string(delivery.TipoEntrega),
	)

	err := s.emailService.SendHTMLEmail(ctx, delivery.Email, subject, htmlBody)
	if err != nil {
		log.Error().
			Err(err).
			Int("delivery_id", delivery.ID).
			Str("email", delivery.Email).
			Msg("Error enviando email de confirmación")
	} else {
		log.Info().
			Int("delivery_id", delivery.ID).
			Str("email", delivery.Email).
			Str("token", delivery.Token).
			Msg("Email de confirmación enviado exitosamente")
	}
}

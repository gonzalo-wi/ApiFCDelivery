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
)

type DeliveryService interface {
	FindAll(ctx context.Context) ([]models.Delivery, error)
	FindByID(ctx context.Context, id int) (*models.Delivery, error)
	FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error)
	FindByRto(ctx context.Context, nroRto string, fechaAccion *time.Time) ([]models.Delivery, error)
	Create(ctx context.Context, delivery *models.Delivery) error
	Update(ctx context.Context, delivery *models.Delivery) error
	Delete(ctx context.Context, id int) error
	CreateFromInfobip(ctx context.Context, req dto.InfobipDeliveryRequest) (*models.Delivery, error)
}
type deliveryService struct {
	store store.DeliveryStore
}

func NewDeliveryService(store store.DeliveryStore) DeliveryService {
	return &deliveryService{store: store}
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

	// Crear dispensers placeholder
	dispensers := createPlaceholderDispensers(req.NroRto, req.Tipos.P, req.Tipos.M)

	// Crear la entrega
	delivery := &models.Delivery{
		NroCta:       req.NroCta,
		NroRto:       req.NroRto,
		Dispensers:   dispensers,
		Cantidad:     cantidadTotal,
		Estado:       models.Pendiente,
		TipoEntrega:  req.TipoEntrega,
		EntregadoPor: req.EntregadoPor,
		SessionID:    req.SessionID,
		FechaAccion:  fechaAccion,
	}

	// Generar token de 4 dígitos (thread-safe)
	delivery.Token = s.generateToken()

	// Guardar en base de datos (la concurrencia se maneja a nivel de BD)
	if err := s.store.Create(ctx, delivery); err != nil {
		return nil, fmt.Errorf("error creando entrega: %w", err)
	}

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

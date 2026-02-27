package service

import (
	"context"
	"fmt"
	"time"

	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"

	"github.com/rs/zerolog/log"
)

type MobileDeliveryService interface {
	ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error)
	ValidateDispenser(ctx context.Context, req dto.ValidateDispenserRequest) (*dto.ValidateDispenserResponse, error)
	CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error)
}

type mobileDeliveryService struct {
	deliveryStore store.DeliveryStore
	publisher     *RabbitMQPublisher
}

func NewMobileDeliveryService(deliveryStore store.DeliveryStore, publisher *RabbitMQPublisher) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore: deliveryStore,
		publisher:     publisher,
	}
}

// ValidateToken valida el token del cliente junto con nro_cta y fecha para mayor seguridad
func (s *mobileDeliveryService) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	// Buscar delivery por token
	deliveries, err := s.deliveryStore.FindAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching deliveries")
		return nil, fmt.Errorf("error validando token: %w", err)
	}

	var foundDelivery *models.Delivery
	for _, d := range deliveries {
		// Validar token + nro_cta + fecha para mayor seguridad
		if d.Token == req.Token &&
			d.NroCta == req.NroCta &&
			d.FechaAccion.Format("2006-01-02") == req.FechaAccion &&
			d.Estado == models.Pendiente {
			foundDelivery = &d
			break
		}
	}

	if foundDelivery == nil {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "Datos de validación incorrectos o entrega ya completada",
		}, nil
	}

	// Construir respuesta
	deliveryInfo := &dto.DeliveryInfoDTO{
		ID:          foundDelivery.ID,
		NroCta:      foundDelivery.NroCta,
		NroRto:      foundDelivery.NroRto,
		Cantidad:    foundDelivery.Cantidad,
		TipoEntrega: string(foundDelivery.TipoEntrega),
		FechaAccion: foundDelivery.FechaAccion.String(),
	}

	dispensers := make([]dto.DispenserInfoDTO, 0, len(foundDelivery.Dispensers))
	for _, d := range foundDelivery.Dispensers {
		dispensers = append(dispensers, dto.DispenserInfoDTO{
			ID:        d.ID,
			Marca:     d.Marca,
			NroSerie:  d.NroSerie,
			Tipo:      string(d.Tipo),
			Validated: false,
		})
	}

	log.Info().
		Int("delivery_id", foundDelivery.ID).
		Str("token", req.Token).
		Str("nro_cta", req.NroCta).
		Msg("Token validated successfully")

	return &dto.ValidateTokenResponse{
		Valid:      true,
		Message:    "Token válido",
		Delivery:   deliveryInfo,
		Dispensers: dispensers,
	}, nil
}

// ValidateDispenser valida que un dispenser escaneado pertenezca al delivery
func (s *mobileDeliveryService) ValidateDispenser(ctx context.Context, req dto.ValidateDispenserRequest) (*dto.ValidateDispenserResponse, error) {
	delivery, err := s.deliveryStore.FindByID(ctx, req.DeliveryID)
	if err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error fetching delivery")
		return nil, fmt.Errorf("error buscando delivery: %w", err)
	}

	// Buscar el dispenser en el delivery
	var found *models.Dispenser
	for _, d := range delivery.Dispensers {
		if d.NroSerie == req.NroSerie {
			found = &d
			break
		}
	}

	if found == nil {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Str("nro_serie", req.NroSerie).
			Msg("Dispenser not found in delivery")

		return &dto.ValidateDispenserResponse{
			Valid:   false,
			Message: "El dispenser no pertenece a esta entrega",
		}, nil
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Str("nro_serie", req.NroSerie).
		Msg("Dispenser validated successfully")

	return &dto.ValidateDispenserResponse{
		Valid:   true,
		Message: "Dispenser válido",
		Dispenser: &dto.DispenserInfoDTO{
			ID:        found.ID,
			Marca:     found.Marca,
			NroSerie:  found.NroSerie,
			Tipo:      string(found.Tipo),
			Validated: true,
		},
	}, nil
}

// CompleteDelivery marca el delivery como completado y publica mensaje a RabbitMQ
func (s *mobileDeliveryService) CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error) {
	// 1. Obtener el delivery
	delivery, err := s.deliveryStore.FindByID(ctx, req.DeliveryID)
	if err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error fetching delivery")
		return nil, fmt.Errorf("error buscando delivery: %w", err)
	}

	// 2. Validar token
	if delivery.Token != req.Token {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Msg("Invalid token for completing delivery")
		return nil, fmt.Errorf("token inválido")
	}

	// 3. Validar estado
	if delivery.Estado != models.Pendiente {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Str("estado", string(delivery.Estado)).
			Msg("Delivery already processed")
		return nil, fmt.Errorf("la entrega ya fue procesada (estado: %s)", delivery.Estado)
	}

	// 4. Validar que todos los dispensers fueron escaneados
	if len(req.Validated) != len(delivery.Dispensers) {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Int("expected", len(delivery.Dispensers)).
			Int("validated", len(req.Validated)).
			Msg("Not all dispensers were validated")
		return nil, fmt.Errorf("faltan dispensers por escanear (esperados: %d, validados: %d)", len(delivery.Dispensers), len(req.Validated))
	}

	// 5. Actualizar estado del delivery
	delivery.Estado = models.Completado
	delivery.UpdatedAt = time.Now()

	if err := s.deliveryStore.Update(ctx, delivery); err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error updating delivery status")
		return nil, fmt.Errorf("error actualizando delivery: %w", err)
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Msg("Delivery marked as completed")

	// 6. Construir mensaje para RabbitMQ
	dispensersMsg := make([]dto.DispenserMessage, 0, len(delivery.Dispensers))
	for _, d := range delivery.Dispensers {
		dispensersMsg = append(dispensersMsg, dto.DispenserMessage{
			Marca:    d.Marca,
			NroSerie: d.NroSerie,
		})
	}

	workOrderMsg := dto.WorkOrderMessageDTO{
		NroCta:     delivery.NroCta,
		Name:       delivery.Name,
		Address:    delivery.Address,
		Locality:   delivery.Locality,
		NroRto:     delivery.NroRto,
		CreatedAt:  delivery.CreatedAt.Format("2006-01-02"),
		TipoAccion: string(delivery.TipoEntrega),
		Token:      delivery.Token,
		Dispensers: dispensersMsg,
		DeliveryID: delivery.ID,
	}

	// 7. Publicar a RabbitMQ
	err = s.publisher.PublishWorkOrder(ctx, workOrderMsg)
	if err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error publishing work order message")
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Msg("Work order message published successfully")

	// 8. Construir respuesta con todos los datos
	dispensersResponse := make([]dto.DispenserCompletedDTO, 0, len(delivery.Dispensers))
	for _, d := range delivery.Dispensers {
		dispensersResponse = append(dispensersResponse, dto.DispenserCompletedDTO{
			Marca:    d.Marca,
			NroSerie: d.NroSerie,
		})
	}

	return &dto.MobileCompleteDeliveryResponse{
		NroCta:     delivery.NroCta,
		Name:       delivery.Name,
		Address:    delivery.Address,
		Locality:   delivery.Locality,
		NroRto:     delivery.NroRto,
		CreatedAt:  delivery.CreatedAt.Format("2006-01-02"),
		TipoAccion: string(delivery.TipoEntrega),
		Token:      delivery.Token,
		Dispensers: dispensersResponse,
	}, nil
}

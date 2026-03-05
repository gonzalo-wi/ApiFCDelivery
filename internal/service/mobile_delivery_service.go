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
	CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error)
	SearchDeliveries(ctx context.Context, fechaAccion string, nroRto string) ([]dto.MobileDeliverySearchResponse, error)
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

// SearchDeliveries busca deliveries por fecha (obligatorio) y opcionalmente por nro_rto
func (s *mobileDeliveryService) SearchDeliveries(ctx context.Context, fechaAccion string, nroRto string) ([]dto.MobileDeliverySearchResponse, error) {
	// Parsear y validar fecha
	parsedDate, err := time.Parse("2006-01-02", fechaAccion)
	if err != nil {
		return nil, fmt.Errorf("fecha inválida. Formato esperado: YYYY-MM-DD")
	}

	// Buscar deliveries
	deliveries, err := s.deliveryStore.FindByRto(ctx, nroRto, &parsedDate)
	if err != nil {
		return nil, fmt.Errorf("error buscando deliveries: %w", err)
	}

	// Mapear a respuesta simplificada
	results := make([]dto.MobileDeliverySearchResponse, 0, len(deliveries))
	for _, d := range deliveries {
		results = append(results, dto.MobileDeliverySearchResponse{
			FechaAccion: d.FechaAccion.Format("2006-01-02"),
			NroCta:      d.NroCta,
			Token:       d.Token,
		})
	}

	return results, nil
}

// ValidateToken valida el token del cliente junto con nro_cta y fecha para mayor seguridad
func (s *mobileDeliveryService) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	// Buscar delivery directamente en BD con índices optimizados
	foundDelivery, err := s.deliveryStore.FindByTokenAndFilters(ctx, req.Token, req.NroCta, req.FechaAccion, models.Pendiente)
	if err != nil {
		log.Error().Err(err).Msg("Error validating token")
		return nil, fmt.Errorf("error validando token: %w", err)
	}

	if foundDelivery == nil {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "Datos de validación incorrectos o entrega ya completada",
		}, nil
	}

	// Construir respuesta con información de items de dispensers
	deliveryInfo := &dto.DeliveryInfoDTO{
		ID:          foundDelivery.ID,
		NroCta:      foundDelivery.NroCta,
		NroRto:      foundDelivery.NroRto,
		Cantidad:    foundDelivery.Cantidad,
		TipoEntrega: string(foundDelivery.TipoEntrega),
		FechaAccion: foundDelivery.FechaAccion.String(),
	}

	itemDispensers := make([]dto.ItemDispenserInfoDTO, 0, len(foundDelivery.ItemDispensers))
	for _, item := range foundDelivery.ItemDispensers {
		itemDispensers = append(itemDispensers, dto.ItemDispenserInfoDTO{
			Tipo:     string(item.Tipo),
			Cantidad: item.Cantidad,
		})
	}

	log.Info().
		Int("delivery_id", foundDelivery.ID).
		Str("token", req.Token).
		Str("nro_cta", req.NroCta).
		Msg("Token validated successfully")

	return &dto.ValidateTokenResponse{
		Valid:          true,
		Message:        "Token válido",
		Delivery:       deliveryInfo,
		ItemDispensers: itemDispensers,
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

	// 4. Actualizar items de dispensers con lo efectivamente entregado
	totalEntregado := uint(0)
	newItemDispensers := make([]models.ItemDispenser, 0, len(req.ItemDispensers))
	for _, item := range req.ItemDispensers {
		tipo := models.TipoDispenser(item.Tipo)
		if tipo != models.TipoDispenserPie && tipo != models.TipoDispenserMesada {
			return nil, fmt.Errorf("tipo de dispenser inválido: %s", item.Tipo)
		}

		newItemDispensers = append(newItemDispensers, models.ItemDispenser{
			Tipo:     tipo,
			Cantidad: item.Cantidad,
		})
		totalEntregado += item.Cantidad
	}

	if totalEntregado == 0 {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Msg("No dispensers were delivered")
		return nil, fmt.Errorf("debe entregar al menos un dispenser")
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Uint("expected", delivery.Cantidad).
		Uint("delivered", totalEntregado).
		Msg("Delivery completed with dispensers")

	// 5. Actualizar datos del cliente desde la app móvil
	if req.Name != "" {
		delivery.Name = req.Name
	}
	if req.Address != "" {
		delivery.Address = req.Address
	}
	if req.Locality != "" {
		delivery.Locality = req.Locality
	}

	// 6. Actualizar items y estado
	delivery.ItemDispensers = newItemDispensers
	delivery.Estado = models.Completado
	delivery.Cantidad = totalEntregado // Cantidad real entregada
	delivery.UpdatedAt = time.Now()

	if err := s.deliveryStore.Update(ctx, delivery); err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error updating delivery status")
		return nil, fmt.Errorf("error actualizando delivery: %w", err)
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Str("name", delivery.Name).
		Str("locality", delivery.Locality).
		Msg("Client data updated from mobile app")

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Msg("Delivery marked as completed")

	// 7. Construir mensaje para RabbitMQ con items de dispensers
	dispensersMsg := make([]dto.DispenserMessage, 0)
	for _, item := range delivery.ItemDispensers {
		// Crear un mensaje por cada dispenser del item
		for i := uint(0); i < item.Cantidad; i++ {
			dispensersMsg = append(dispensersMsg, dto.DispenserMessage{
				Marca:    fmt.Sprintf("Dispenser-%s", item.Tipo),
				NroSerie: fmt.Sprintf("%s-%d", item.Tipo, i+1),
			})
		}
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

	// 8. Publicar a RabbitMQ
	workOrderQueued := false
	err = s.publisher.PublishWorkOrder(ctx, workOrderMsg)
	if err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error publishing work order message")
	} else {
		workOrderQueued = true
		log.Info().
			Int("delivery_id", req.DeliveryID).
			Msg("Work order message published successfully")
	}

	// 9. Construir respuesta con items de dispensers entregados
	itemDispensersResponse := make([]dto.ItemDispenserCompletedDTO, 0, len(delivery.ItemDispensers))
	for _, item := range delivery.ItemDispensers {
		itemDispensersResponse = append(itemDispensersResponse, dto.ItemDispenserCompletedDTO{
			Tipo:     string(item.Tipo),
			Cantidad: item.Cantidad,
		})
	}

	return &dto.MobileCompleteDeliveryResponse{
		NroCta:          delivery.NroCta,
		Name:            delivery.Name,
		Address:         delivery.Address,
		Locality:        delivery.Locality,
		NroRto:          delivery.NroRto,
		CreatedAt:       delivery.CreatedAt.Format("2006-01-02"),
		TipoAccion:      string(delivery.TipoEntrega),
		Token:           delivery.Token,
		ItemDispensers:  itemDispensersResponse,
		WorkOrderQueued: workOrderQueued,
	}, nil
}

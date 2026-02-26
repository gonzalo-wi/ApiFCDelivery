package service

import (
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type DeliveryWithTermsService interface {
	InitiateDelivery(ctx context.Context, req dto.InitiateDeliveryRequest, appBaseURL string, ttlHours int) (*dto.InitiateDeliveryResponse, error)
	CompleteDelivery(ctx context.Context, termsToken string) (*models.Delivery, error)
	GetDeliveryByTermsToken(ctx context.Context, termsToken string) (*models.Delivery, error)
}

type deliveryWithTermsService struct {
	deliveryStore       store.DeliveryStore
	termsSessionStore   store.TermsSessionStore
	termsSessionService TermsSessionService
}

func NewDeliveryWithTermsService(
	deliveryStore store.DeliveryStore,
	termsSessionStore store.TermsSessionStore,
	termsSessionService TermsSessionService,
) DeliveryWithTermsService {
	return &deliveryWithTermsService{
		deliveryStore:       deliveryStore,
		termsSessionStore:   termsSessionStore,
		termsSessionService: termsSessionService,
	}
}

// InitiateDelivery crea una sesión de términos y guarda los datos de la entrega pendiente
func (s *deliveryWithTermsService) InitiateDelivery(
	ctx context.Context,
	req dto.InitiateDeliveryRequest,
	appBaseURL string,
	ttlHours int,
) (*dto.InitiateDeliveryResponse, error) {

	deliveryData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error serializando datos de entrega: %w", err)
	}

	sessionID := req.NroRto

	termsResponse, err := s.termsSessionService.CreateSession(ctx, sessionID, appBaseURL, ttlHours)
	if err != nil {
		return nil, fmt.Errorf("error creando sesión de términos: %w", err)
	}

	termsSession, err := s.termsSessionStore.FindByToken(ctx, termsResponse.Token)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo sesión de términos: %w", err)
	}

	termsSession.DeliveryData = string(deliveryData)
	if err := s.termsSessionStore.Update(ctx, termsSession); err != nil {
		return nil, fmt.Errorf("error actualizando sesión con datos de entrega: %w", err)
	}

	log.Info().
		Str("token", termsResponse.Token).
		Str("nro_rto", req.NroRto).
		Msg("Entrega iniciada - esperando aceptación de términos")

	return &dto.InitiateDeliveryResponse{
		Token:     termsResponse.Token,
		TermsURL:  termsResponse.URL,
		ExpiresAt: termsResponse.ExpiresAt.Format(time.RFC3339),
		Message:   "Por favor, acepte los términos y condiciones para completar la entrega",
	}, nil
}

// CompleteDelivery crea la entrega después de que los términos fueron aceptados
func (s *deliveryWithTermsService) CompleteDelivery(ctx context.Context, termsToken string) (*models.Delivery, error) {

	termsSession, err := s.termsSessionStore.FindByToken(ctx, termsToken)
	if err != nil {
		return nil, fmt.Errorf("sesión de términos no encontrada")
	}

	if termsSession.Status != models.StatusAccepted {
		return nil, fmt.Errorf("los términos no han sido aceptados (estado: %s)", termsSession.Status)
	}

	if time.Now().After(termsSession.ExpiresAt) {
		return nil, fmt.Errorf("la sesión de términos ha expirado")
	}

	if termsSession.DeliveryData == "" {
		return nil, fmt.Errorf("no hay datos de entrega asociados a esta sesión")
	}

	var deliveryReq dto.InitiateDeliveryRequest
	if err := json.Unmarshal([]byte(termsSession.DeliveryData), &deliveryReq); err != nil {
		return nil, fmt.Errorf("error deserializando datos de entrega: %w", err)
	}

	// Parsear fecha de acción usando helper
	fechaAccion, err := parseFechaAccion(deliveryReq.FechaAccion)
	if err != nil {
		return nil, err
	}

	// Crear dispensers
	dispensers := make([]models.Dispenser, len(deliveryReq.Dispensers))
	for i, d := range deliveryReq.Dispensers {
		dispensers[i] = models.Dispenser{
			Marca:    d.Marca,
			NroSerie: d.NroSerie,
			Tipo:     d.Tipo,
		}
	}

	// Crear la entrega
	delivery := &models.Delivery{
		NroCta:         deliveryReq.NroCta,
		NroRto:         deliveryReq.NroRto,
		Dispensers:     dispensers,
		Cantidad:       deliveryReq.Cantidad,
		Estado:         models.Completado,
		TipoEntrega:    deliveryReq.TipoEntrega,
		TermsSessionID: &termsSession.ID,
		FechaAccion:    fechaAccion,
	}

	// Generar token de 4 dígitos (el que ya existe en delivery service)
	delivery.Token = generateDeliveryToken()

	// Guardar en BD
	if err := s.deliveryStore.Create(ctx, delivery); err != nil {
		return nil, fmt.Errorf("error creando entrega: %w", err)
	}

	log.Info().
		Int("delivery_id", delivery.ID).
		Str("nro_rto", delivery.NroRto).
		Str("terms_token", termsToken).
		Msg("Entrega completada exitosamente después de aceptar términos")

	return delivery, nil
}

// GetDeliveryByTermsToken obtiene la entrega asociada a un token de términos
func (s *deliveryWithTermsService) GetDeliveryByTermsToken(ctx context.Context, termsToken string) (*models.Delivery, error) {
	termsSession, err := s.termsSessionStore.FindByToken(ctx, termsToken)
	if err != nil {
		return nil, fmt.Errorf("sesión de términos no encontrada")
	}
	if termsSession.DeliveryData == "" {
		return nil, fmt.Errorf("no hay datos de entrega asociados")
	}
	return nil, fmt.Errorf("funcionalidad no implementada")
}

// Función auxiliar para generar token de entrega (4 dígitos)
func generateDeliveryToken() string {
	return fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
}

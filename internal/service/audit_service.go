package service

import (
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AuditService struct {
	auditStore *store.AuditEventStore
}

func NewAuditService(auditStore *store.AuditEventStore) *AuditService {
	return &AuditService{
		auditStore: auditStore,
	}
}

// LogEvent registra un evento de auditoría de forma síncrona
func (s *AuditService) LogEvent(ctx context.Context, event *models.AuditEvent) error {
	return s.auditStore.Create(ctx, event)
}

// LogEventAsync registra un evento de auditoría de forma asíncrona (no bloqueante)
// Usar cuando la auditoría no debe afectar el tiempo de respuesta
func (s *AuditService) LogEventAsync(event *models.AuditEvent) {
	s.auditStore.CreateAsync(event)
}

// GetEntityHistory obtiene el historial de cambios de una entidad
func (s *AuditService) GetEntityHistory(ctx context.Context, entityType models.AuditEntityType, entityID string, limit int) ([]*models.AuditEvent, error) {
	return s.auditStore.FindByEntity(ctx, string(entityType), entityID, limit)
}

// GetActorActivity obtiene la actividad de un actor
func (s *AuditService) GetActorActivity(ctx context.Context, actorType models.AuditActorType, actorID string, limit int) ([]*models.AuditEvent, error) {
	return s.auditStore.FindByActor(ctx, string(actorType), actorID, limit)
}

// GetRequestTrace obtiene todos los eventos asociados a un request
func (s *AuditService) GetRequestTrace(ctx context.Context, requestID uuid.UUID) ([]*models.AuditEvent, error) {
	return s.auditStore.FindByRequestID(ctx, requestID)
}

// GetRecentEvents obtiene eventos recientes
func (s *AuditService) GetRecentEvents(ctx context.Context, hours int, limit int) ([]*models.AuditEvent, error) {
	return s.auditStore.FindRecent(ctx, hours, limit)
}

// Search realiza búsqueda avanzada de eventos
func (s *AuditService) Search(ctx context.Context, filter store.AuditSearchFilter) ([]*models.AuditEvent, int, error) {
	return s.auditStore.Search(ctx, filter)
}

// CleanupOldEvents limpia eventos antiguos según política de retención
func (s *AuditService) CleanupOldEvents(ctx context.Context, retentionMonths int) (int64, error) {
	return s.auditStore.CleanupOld(ctx, retentionMonths)
}

// Helper functions para crear eventos comunes

// LogDeliveryCreated registra la creación de una entrega
func (s *AuditService) LogDeliveryCreated(ctx context.Context, deliveryID int, actorType models.AuditActorType, actorID string, delivery interface{}, metadata map[string]interface{}) {
	event := models.NewAuditEvent().
		WithEntity(models.EntityDelivery, fmt.Sprintf("%d", deliveryID)).
		WithAction(models.ActionCreated).
		WithActor(actorType, actorID).
		WithAfterState(delivery).
		WithMetadata(metadata).
		Build()

	s.LogEventAsync(event)
}

// LogDeliveryUpdated registra la actualización de una entrega
func (s *AuditService) LogDeliveryUpdated(ctx context.Context, deliveryID int, actorType models.AuditActorType, actorID string, before, after interface{}, metadata map[string]interface{}) {
	event := models.NewAuditEvent().
		WithEntity(models.EntityDelivery, fmt.Sprintf("%d", deliveryID)).
		WithAction(models.ActionUpdated).
		WithActor(actorType, actorID).
		WithBeforeState(before).
		WithAfterState(after).
		WithMetadata(metadata).
		Build()

	s.LogEventAsync(event)
}

// LogTokenGenerated registra la generación de un token
func (s *AuditService) LogTokenGenerated(ctx context.Context, provider string, ipAddress, userAgent string) {
	event := models.NewAuditEvent().
		WithEntity(models.EntityAuthToken, provider).
		WithAction(models.ActionCreated).
		WithActor(models.ActorExternal, provider).
		WithHTTPContext(ipAddress, userAgent).
		WithMetadata(map[string]interface{}{
			"provider": provider,
		}).
		Build()

	s.LogEventAsync(event)
}

// LogTokenValidated registra la validación de un token
func (s *AuditService) LogTokenValidated(ctx context.Context, provider string, valid bool, ipAddress, path string) {
	action := models.ActionValidated
	if !valid {
		action = models.ActionRejected
	}

	event := models.NewAuditEvent().
		WithEntity(models.EntityAuthToken, provider).
		WithAction(action).
		WithActor(models.ActorExternal, provider).
		WithHTTPContext(ipAddress, "").
		WithMetadata(map[string]interface{}{
			"valid":    valid,
			"endpoint": path,
		}).
		Build()

	s.LogEventAsync(event)
}

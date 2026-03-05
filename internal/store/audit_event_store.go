package store

import (
	"GoFrioCalor/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type AuditEventStore struct {
	db *sqlx.DB
}

func NewAuditEventStore(db *sqlx.DB) *AuditEventStore {
	return &AuditEventStore{db: db}
}

// Create inserta un nuevo evento de auditoría
func (s *AuditEventStore) Create(ctx context.Context, event *models.AuditEvent) error {
	query := `
		INSERT INTO audit_events (
			id, occurred_at, service, entity_type, entity_id, action,
			actor_type, actor_id, request_id, trace_id, ip_address, user_agent,
			before_state, after_state, metadata, created_at
		) VALUES (
			:id, :occurred_at, :service, :entity_type, :entity_id, :action,
			:actor_type, :actor_id, :request_id, :trace_id, :ip_address, :user_agent,
			:before_state, :after_state, :metadata, :created_at
		)`

	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		log.Error().
			Err(err).
			Str("entity_type", event.EntityType).
			Str("entity_id", event.EntityID).
			Str("action", event.Action).
			Msg("Error creating audit event")
		return err
	}

	return nil
}

// CreateAsync inserta un evento de auditoría de forma asíncrona (no bloqueante)
func (s *AuditEventStore) CreateAsync(event *models.AuditEvent) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Create(ctx, event); err != nil {
			log.Error().
				Err(err).
				Str("entity_type", event.EntityType).
				Str("action", event.Action).
				Msg("Failed to create audit event asynchronously")
		}
	}()
}

// FindByEntity obtiene eventos de auditoría para una entidad específica
func (s *AuditEventStore) FindByEntity(ctx context.Context, entityType, entityID string, limit int) ([]*models.AuditEvent, error) {
	query := `
		SELECT * FROM audit_events
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY occurred_at DESC
		LIMIT $3`

	var events []*models.AuditEvent
	err := s.db.SelectContext(ctx, &events, query, entityType, entityID, limit)
	if err != nil {
		log.Error().
			Err(err).
			Str("entity_type", entityType).
			Str("entity_id", entityID).
			Msg("Error finding audit events by entity")
		return nil, err
	}

	return events, nil
}

// FindByActor obtiene eventos de auditoría para un actor específico
func (s *AuditEventStore) FindByActor(ctx context.Context, actorType, actorID string, limit int) ([]*models.AuditEvent, error) {
	query := `
		SELECT * FROM audit_events
		WHERE actor_type = $1 AND actor_id = $2
		ORDER BY occurred_at DESC
		LIMIT $3`

	var events []*models.AuditEvent
	err := s.db.SelectContext(ctx, &events, query, actorType, actorID, limit)
	if err != nil {
		log.Error().
			Err(err).
			Str("actor_type", actorType).
			Str("actor_id", actorID).
			Msg("Error finding audit events by actor")
		return nil, err
	}

	return events, nil
}

// FindByRequestID obtiene todos los eventos asociados a un request ID
func (s *AuditEventStore) FindByRequestID(ctx context.Context, requestID uuid.UUID) ([]*models.AuditEvent, error) {
	query := `
		SELECT * FROM audit_events
		WHERE request_id = $1
		ORDER BY occurred_at ASC`

	var events []*models.AuditEvent
	err := s.db.SelectContext(ctx, &events, query, requestID)
	if err != nil {
		log.Error().
			Err(err).
			Str("request_id", requestID.String()).
			Msg("Error finding audit events by request ID")
		return nil, err
	}

	return events, nil
}

// FindRecent obtiene eventos recientes (últimas N horas)
func (s *AuditEventStore) FindRecent(ctx context.Context, hours int, limit int) ([]*models.AuditEvent, error) {
	query := `
		SELECT * FROM audit_events
		WHERE occurred_at > NOW() - ($1 || ' hours')::INTERVAL
		ORDER BY occurred_at DESC
		LIMIT $2`

	var events []*models.AuditEvent
	err := s.db.SelectContext(ctx, &events, query, hours, limit)
	if err != nil {
		log.Error().
			Err(err).
			Int("hours", hours).
			Msg("Error finding recent audit events")
		return nil, err
	}

	return events, nil
}

// AuditSearchFilter filtros para búsqueda avanzada
type AuditSearchFilter struct {
	EntityType *string
	EntityID   *string
	Action     *string
	ActorType  *string
	ActorID    *string
	FromDate   *time.Time
	ToDate     *time.Time
	Limit      int
	Offset     int
}

// Search búsqueda avanzada de eventos de auditoría
func (s *AuditEventStore) Search(ctx context.Context, filter AuditSearchFilter) ([]*models.AuditEvent, int, error) {
	// Query base
	baseQuery := "FROM audit_events WHERE 1=1"
	args := make(map[string]interface{})

	// Construir WHERE dinámicamente
	whereClause := ""

	if filter.EntityType != nil {
		whereClause += " AND entity_type = :entity_type"
		args["entity_type"] = *filter.EntityType
	}

	if filter.EntityID != nil {
		whereClause += " AND entity_id = :entity_id"
		args["entity_id"] = *filter.EntityID
	}

	if filter.Action != nil {
		whereClause += " AND action = :action"
		args["action"] = *filter.Action
	}

	if filter.ActorType != nil {
		whereClause += " AND actor_type = :actor_type"
		args["actor_type"] = *filter.ActorType
	}

	if filter.ActorID != nil {
		whereClause += " AND actor_id = :actor_id"
		args["actor_id"] = *filter.ActorID
	}

	if filter.FromDate != nil {
		whereClause += " AND occurred_at >= :from_date"
		args["from_date"] = *filter.FromDate
	}

	if filter.ToDate != nil {
		whereClause += " AND occurred_at <= :to_date"
		args["to_date"] = *filter.ToDate
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var total int

	countStmt, err := s.db.PrepareNamedContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countStmt.Close()

	err = countStmt.GetContext(ctx, &total, args)
	if err != nil {
		log.Error().Err(err).Msg("Error counting audit events")
		return nil, 0, err
	}

	// Select query
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	args["limit"] = filter.Limit
	args["offset"] = filter.Offset

	selectQuery := "SELECT * " + baseQuery + whereClause + " ORDER BY occurred_at DESC LIMIT :limit OFFSET :offset"

	stmt, err := s.db.PrepareNamedContext(ctx, selectQuery)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	var events []*models.AuditEvent
	err = stmt.SelectContext(ctx, &events, args)
	if err != nil {
		log.Error().Err(err).Msg("Error searching audit events")
		return nil, 0, err
	}

	return events, total, nil
}

// CleanupOld elimina eventos antiguos según política de retención
func (s *AuditEventStore) CleanupOld(ctx context.Context, retentionMonths int) (int64, error) {
	query := `SELECT cleanup_old_audit_events($1)`

	var deletedCount int64
	err := s.db.GetContext(ctx, &deletedCount, query, retentionMonths)
	if err != nil {
		log.Error().
			Err(err).
			Int("retention_months", retentionMonths).
			Msg("Error cleaning up old audit events")
		return 0, err
	}

	log.Info().
		Int64("deleted_count", deletedCount).
		Int("retention_months", retentionMonths).
		Msg("Cleanup old audit events completed")

	return deletedCount, nil
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditEvent representa un evento de auditoría en el sistema
type AuditEvent struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	OccurredAt  time.Time  `json:"occurred_at" db:"occurred_at"`
	Service     string     `json:"service" db:"service"`
	EntityType  string     `json:"entity_type" db:"entity_type"`
	EntityID    string     `json:"entity_id" db:"entity_id"`
	Action      string     `json:"action" db:"action"`
	ActorType   *string    `json:"actor_type,omitempty" db:"actor_type"`
	ActorID     *string    `json:"actor_id,omitempty" db:"actor_id"`
	RequestID   *uuid.UUID `json:"request_id,omitempty" db:"request_id"`
	TraceID     *string    `json:"trace_id,omitempty" db:"trace_id"`
	IPAddress   *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent   *string    `json:"user_agent,omitempty" db:"user_agent"`
	BeforeState JSONB      `json:"before_state,omitempty" db:"before_state"`
	AfterState  JSONB      `json:"after_state,omitempty" db:"after_state"`
	Metadata    JSONB      `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// JSONB tipo personalizado para campos JSONB de PostgreSQL
type JSONB map[string]interface{}

// Value implementa driver.Valuer para JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implementa sql.Scanner para JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(JSONB)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

// AuditAction define las acciones auditables
type AuditAction string

const (
	ActionCreated   AuditAction = "CREATED"
	ActionUpdated   AuditAction = "UPDATED"
	ActionDeleted   AuditAction = "DELETED"
	ActionViewed    AuditAction = "VIEWED"
	ActionAssigned  AuditAction = "ASSIGNED"
	ActionDelivered AuditAction = "DELIVERED"
	ActionCanceled  AuditAction = "CANCELED"
	ActionValidated AuditAction = "VALIDATED"
	ActionSent      AuditAction = "SENT"
	ActionAccepted  AuditAction = "ACCEPTED"
	ActionRejected  AuditAction = "REJECTED"
)

// AuditEntityType define los tipos de entidades auditables
type AuditEntityType string

const (
	EntityDelivery       AuditEntityType = "delivery"
	EntityWorkOrder      AuditEntityType = "work_order"
	EntityTermsSession   AuditEntityType = "terms_session"
	EntityAuthToken      AuditEntityType = "auth_token"
	EntityItemDispenser  AuditEntityType = "item_dispenser"
	EntityTruck          AuditEntityType = "truck"
	EntityMobileDelivery AuditEntityType = "mobile_delivery"
)

// AuditActorType define los tipos de actores
type AuditActorType string

const (
	ActorSystem    AuditActorType = "system"
	ActorAPIClient AuditActorType = "api_client"
	ActorMobileApp AuditActorType = "mobile_app"
	ActorAdmin     AuditActorType = "admin"
	ActorReparto   AuditActorType = "reparto"
	ActorExternal  AuditActorType = "external"
)

// AuditEventBuilder facilita la creación de eventos de auditoría
type AuditEventBuilder struct {
	event *AuditEvent
}

// NewAuditEvent crea un nuevo builder para eventos de auditoría
func NewAuditEvent() *AuditEventBuilder {
	return &AuditEventBuilder{
		event: &AuditEvent{
			ID:         uuid.New(),
			OccurredAt: time.Now(),
			Service:    "dispenser-api",
			CreatedAt:  time.Now(),
		},
	}
}

// WithEntity establece el tipo y ID de entidad
func (b *AuditEventBuilder) WithEntity(entityType AuditEntityType, entityID string) *AuditEventBuilder {
	b.event.EntityType = string(entityType)
	b.event.EntityID = entityID
	return b
}

// WithAction establece la acción
func (b *AuditEventBuilder) WithAction(action AuditAction) *AuditEventBuilder {
	b.event.Action = string(action)
	return b
}

// WithActor establece el tipo y ID del actor
func (b *AuditEventBuilder) WithActor(actorType AuditActorType, actorID string) *AuditEventBuilder {
	actorTypeStr := string(actorType)
	b.event.ActorType = &actorTypeStr
	b.event.ActorID = &actorID
	return b
}

// WithRequestID establece el ID del request
func (b *AuditEventBuilder) WithRequestID(requestID uuid.UUID) *AuditEventBuilder {
	b.event.RequestID = &requestID
	return b
}

// WithTraceID establece el trace ID
func (b *AuditEventBuilder) WithTraceID(traceID string) *AuditEventBuilder {
	b.event.TraceID = &traceID
	return b
}

// WithHTTPContext establece IP y user agent
func (b *AuditEventBuilder) WithHTTPContext(ip, userAgent string) *AuditEventBuilder {
	if ip != "" {
		b.event.IPAddress = &ip
	}
	if userAgent != "" {
		b.event.UserAgent = &userAgent
	}
	return b
}

// WithBeforeState establece el estado anterior
func (b *AuditEventBuilder) WithBeforeState(state interface{}) *AuditEventBuilder {
	if state != nil {
		b.event.BeforeState = toJSONB(state)
	}
	return b
}

// WithAfterState establece el estado posterior
func (b *AuditEventBuilder) WithAfterState(state interface{}) *AuditEventBuilder {
	if state != nil {
		b.event.AfterState = toJSONB(state)
	}
	return b
}

// WithMetadata establece metadata adicional
func (b *AuditEventBuilder) WithMetadata(metadata map[string]interface{}) *AuditEventBuilder {
	if metadata != nil {
		b.event.Metadata = JSONB(metadata)
	}
	return b
}

// Build retorna el evento de auditoría construido
func (b *AuditEventBuilder) Build() *AuditEvent {
	return b.event
}

// toJSONB convierte cualquier estructura a JSONB
func toJSONB(v interface{}) JSONB {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	var result JSONB
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil
	}

	return result
}

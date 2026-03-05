-- Migration: Create audit_events table
-- Description: Sistema de auditoría completo para tracking de todas las operaciones

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabla principal de auditoría
CREATE TABLE IF NOT EXISTS audit_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Contexto del servicio
    service VARCHAR(50) NOT NULL DEFAULT 'gofrioCalor',
    
    -- Entidad afectada
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    
    -- Acción realizada
    action VARCHAR(50) NOT NULL,
    
    -- Actor que realizó la acción
    actor_type VARCHAR(50),
    actor_id VARCHAR(100),
    
    -- Trazabilidad
    request_id UUID,
    trace_id VARCHAR(100),
    
    -- Contexto HTTP (cuando aplica)
    ip_address INET,
    user_agent TEXT,
    
    -- Estados antes/después (para auditoría completa)
    before_state JSONB,
    after_state JSONB,
    
    -- Metadata adicional
    metadata JSONB,
    
    -- Índices para búsquedas comunes
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índices para optimizar queries comunes
CREATE INDEX idx_audit_events_occurred_at ON audit_events (occurred_at DESC);
CREATE INDEX idx_audit_events_entity ON audit_events (entity_type, entity_id);
CREATE INDEX idx_audit_events_actor ON audit_events (actor_type, actor_id);
CREATE INDEX idx_audit_events_action ON audit_events (action);
CREATE INDEX idx_audit_events_request_id ON audit_events (request_id) WHERE request_id IS NOT NULL;
CREATE INDEX idx_audit_events_service ON audit_events (service);

-- Índice compuesto para búsquedas por fecha y tipo
CREATE INDEX idx_audit_events_occurred_entity_type ON audit_events (occurred_at DESC, entity_type);

-- Índice JSONB para búsquedas en metadata
CREATE INDEX idx_audit_events_metadata ON audit_events USING gin (metadata);

-- Comentarios para documentación
COMMENT ON TABLE audit_events IS 'Registro completo de auditoría para todas las operaciones del sistema';
COMMENT ON COLUMN audit_events.id IS 'Identificador único del evento de auditoría';
COMMENT ON COLUMN audit_events.occurred_at IS 'Timestamp con zona horaria de cuando ocurrió el evento';
COMMENT ON COLUMN audit_events.service IS 'Nombre del servicio que generó el evento';
COMMENT ON COLUMN audit_events.entity_type IS 'Tipo de entidad (Delivery, WorkOrder, TermsSession, etc.)';
COMMENT ON COLUMN audit_events.entity_id IS 'ID de la entidad afectada';
COMMENT ON COLUMN audit_events.action IS 'Acción realizada (CREATED, UPDATED, DELETED, etc.)';
COMMENT ON COLUMN audit_events.actor_type IS 'Tipo de actor (user, system, api, mobile_app, etc.)';
COMMENT ON COLUMN audit_events.actor_id IS 'ID del actor que realizó la acción';
COMMENT ON COLUMN audit_events.request_id IS 'ID único del request HTTP para correlación';
COMMENT ON COLUMN audit_events.trace_id IS 'ID de traza distribuida para microservicios';
COMMENT ON COLUMN audit_events.ip_address IS 'Dirección IP del cliente (cuando aplica)';
COMMENT ON COLUMN audit_events.user_agent IS 'User agent del cliente (cuando aplica)';
COMMENT ON COLUMN audit_events.before_state IS 'Estado de la entidad antes del cambio (JSONB)';
COMMENT ON COLUMN audit_events.after_state IS 'Estado de la entidad después del cambio (JSONB)';
COMMENT ON COLUMN audit_events.metadata IS 'Información adicional contextual (motivos, cantidades, etc.)';

-- Opcional: Particionamiento por rango temporal (para grandes volúmenes)
-- Descomenta si esperas >1M eventos/mes
/*
ALTER TABLE audit_events RENAME TO audit_events_template;

CREATE TABLE audit_events (
    LIKE audit_events_template INCLUDING ALL
) PARTITION BY RANGE (occurred_at);

-- Crear particiones mensuales
CREATE TABLE audit_events_2026_03 PARTITION OF audit_events
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE TABLE audit_events_2026_04 PARTITION OF audit_events
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
*/

-- Función para limpiar eventos antiguos (retención de 2 años por defecto)
CREATE OR REPLACE FUNCTION cleanup_old_audit_events(retention_months INT DEFAULT 24)
RETURNS INT AS $$
DECLARE
    deleted_count INT;
BEGIN
    DELETE FROM audit_events
    WHERE occurred_at < NOW() - (retention_months || ' months')::INTERVAL;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_old_audit_events IS 'Limpia eventos de auditoría antiguos según política de retención';

-- Vista para eventos recientes (últimas 24 horas)
CREATE OR REPLACE VIEW recent_audit_events AS
SELECT *
FROM audit_events
WHERE occurred_at > NOW() - INTERVAL '24 hours'
ORDER BY occurred_at DESC;

COMMENT ON VIEW recent_audit_events IS 'Vista con eventos de auditoría de las últimas 24 horas';

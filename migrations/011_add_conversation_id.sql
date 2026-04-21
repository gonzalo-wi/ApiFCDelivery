-- Migración: Reemplazar session_id por conversation_id como clave de idempotencia
-- Contexto: Infobip confirmó que session_id NO es único entre conversaciones.
--           El campo conversation_id SÍ es único y es el identificador correcto.
--
-- Cambios en deliveries:
--   - Renombrar session_id → conversation_id
--   - Renombrar el índice correspondiente
--
-- Cambios en terms_sessions:
--   - Agregar conversation_id (nullable, único donde no es NULL)
--   - session_id se mantiene para el callback webhook de Infobip

-- ========== deliveries ==========

-- 1. Renombrar la columna
ALTER TABLE deliveries RENAME COLUMN session_id TO conversation_id;

-- 2. Renombrar el índice único
ALTER INDEX IF EXISTS idx_session_id RENAME TO idx_conversation_id;

-- Si el índice tenía otro nombre por GORM, usar este fallback:
-- DROP INDEX IF EXISTS idx_deliveries_session_id;
-- CREATE UNIQUE INDEX idx_conversation_id ON deliveries(conversation_id)
--   WHERE conversation_id IS NOT NULL AND conversation_id != '';

-- ========== terms_sessions ==========

-- 3. Agregar la columna conversation_id (nullable)
ALTER TABLE terms_sessions ADD COLUMN IF NOT EXISTS conversation_id VARCHAR(255) NULL;

-- 4. Índice único sobre conversation_id (ignorar NULLs — comportamiento default en PostgreSQL)
CREATE UNIQUE INDEX IF NOT EXISTS idx_terms_conversation_id
    ON terms_sessions(conversation_id)
    WHERE conversation_id IS NOT NULL;

-- Verificación (ejecutar manualmente si se desea):
-- SELECT conversation_id, COUNT(*) FROM deliveries
--   WHERE conversation_id IS NOT NULL GROUP BY conversation_id HAVING COUNT(*) > 1;
-- SELECT conversation_id, COUNT(*) FROM terms_sessions
--   WHERE conversation_id IS NOT NULL GROUP BY conversation_id HAVING COUNT(*) > 1;

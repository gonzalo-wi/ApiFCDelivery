-- Migración: Agregar índice único a session_id para prevenir entregas duplicadas desde Infobip
-- Fecha: 2026-02-27

-- Eliminar el índice regular existente si existe (ignora error si no existe)
ALTER TABLE deliveries DROP INDEX idx_deliveries_session_id;

-- Crear índice único en session_id para prevenir duplicados
-- Nota: En MySQL, los índices únicos ignoran múltiples valores NULL automáticamente
CREATE UNIQUE INDEX idx_session_id ON deliveries(session_id);

-- Verificación (descomentar para ejecutar manualmente)
-- SELECT session_id, COUNT(*) as count 
-- FROM deliveries 
-- WHERE session_id IS NOT NULL AND session_id != ''
-- GROUP BY session_id 
-- HAVING COUNT(*) > 1;

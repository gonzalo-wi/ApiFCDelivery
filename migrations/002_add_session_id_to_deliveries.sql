-- Agregar columna session_id a la tabla deliveries para relacionar con Infobip
-- Este campo permite machear la sesión de aceptación de términos con la entrega

ALTER TABLE deliveries 
ADD COLUMN session_id VARCHAR(255) NULL AFTER tipo_entrega,
ADD INDEX idx_session_id (session_id);

-- Nota: El campo es opcional (NULL) porque deliveries existentes no tienen session_id
-- y no todas las entregas requieren aceptación de términos

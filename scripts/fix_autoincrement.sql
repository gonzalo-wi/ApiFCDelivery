-- Script para reparar el AUTO_INCREMENT de terms_sessions
-- Ejecutar este script en MySQL/MariaDB

USE friocalor;

-- Ver el estado actual de la tabla
SHOW CREATE TABLE terms_sessions;

-- Reparar AUTO_INCREMENT
ALTER TABLE terms_sessions AUTO_INCREMENT = 1;

-- Cambiar el tipo de ID a BIGINT para evitar problemas futuros
ALTER TABLE terms_sessions MODIFY id BIGINT AUTO_INCREMENT;

-- Verificar que el cambio se aplicó correctamente
SHOW CREATE TABLE terms_sessions;

-- Ver registros existentes
SELECT COUNT(*) as total_registros FROM terms_sessions;

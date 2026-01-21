-- Tabla para gestionar sesiones de términos y condiciones con Infobip
CREATE TABLE IF NOT EXISTS terms_sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(64) NOT NULL UNIQUE,
    session_id VARCHAR(255) NOT NULL,
    status ENUM('PENDING', 'ACCEPTED', 'REJECTED', 'EXPIRED') NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP NULL,
    rejected_at TIMESTAMP NULL,
    ip VARCHAR(45) NULL,
    user_agent TEXT NULL,
    notify_status ENUM('PENDING', 'SENT', 'FAILED') NOT NULL DEFAULT 'PENDING',
    notify_attempts INT NOT NULL DEFAULT 0,
    last_error TEXT NULL,
    
    INDEX idx_token (token),
    INDEX idx_session_id (session_id),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Nota: Si usas GORM con AutoMigrate, esta tabla se creará automáticamente.
-- Este script SQL es solo para referencia o para crear manualmente la tabla.

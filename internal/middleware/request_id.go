package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

// RequestID middleware agrega un ID único a cada request para trazabilidad
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Intentar obtener el request ID del header (si el cliente lo envía)
		requestID := c.GetHeader(RequestIDHeader)

		// Si no existe, generar uno nuevo
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Guardar en el context de Gin
		c.Set(RequestIDKey, requestID)

		// Agregar al response header para que el cliente pueda correlacionar
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID obtiene el request ID del context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

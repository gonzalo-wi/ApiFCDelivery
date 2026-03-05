package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

// RegisterAuditRoutes registra las rutas de auditoría
func RegisterAuditRoutes(api *gin.RouterGroup, auditHandler *transport.AuditHandler) {
	audit := api.Group("/audit")
	{
		// Historial de una entidad específica
		audit.GET("/entity/:entity_type/:entity_id", auditHandler.GetEntityHistory)

		// Actividad de un actor específico
		audit.GET("/actor/:actor_type/:actor_id", auditHandler.GetActorActivity)

		// Traza completa de un request
		audit.GET("/request/:request_id", auditHandler.GetRequestTrace)

		// Eventos recientes
		audit.GET("/recent", auditHandler.GetRecentEvents)

		// Búsqueda avanzada
		audit.POST("/search", auditHandler.Search)

		// Estadísticas
		audit.GET("/stats", auditHandler.GetStats)
	}
}

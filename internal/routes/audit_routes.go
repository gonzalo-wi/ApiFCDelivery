package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterAuditRoutes(api *gin.RouterGroup, auditHandler *transport.AuditHandler) {
	audit := api.Group("/audit")
	{
		audit.GET("/entity/:entity_type/:entity_id", auditHandler.GetEntityHistory)
		audit.GET("/actor/:actor_type/:actor_id", auditHandler.GetActorActivity)
		audit.GET("/request/:request_id", auditHandler.GetRequestTrace)
		audit.GET("/recent", auditHandler.GetRecentEvents)
		audit.POST("/search", auditHandler.Search)
		audit.GET("/stats", auditHandler.GetStats)
	}
}

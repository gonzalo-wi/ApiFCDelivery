package transport

import (
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/service"
	"GoFrioCalor/internal/store"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

// GetEntityHistory obtiene el historial de auditoría de una entidad
// GET /audit/entity/:entity_type/:entity_id
func (h *AuditHandler) GetEntityHistory(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityID := c.Param("entity_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	} else if limit > 500 {
		limit = 500 // Máximo permitido
	}

	events, err := h.auditService.GetEntityHistory(
		c.Request.Context(),
		models.AuditEntityType(entityType),
		entityID,
		limit,
	)

	if err != nil {
		log.Error().Err(err).Msg("Error getting entity history")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo historial de auditoría",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entity_type": entityType,
		"entity_id":   entityID,
		"total":       len(events),
		"events":      events,
	})
}

// GetActorActivity obtiene la actividad de un actor
// GET /audit/actor/:actor_type/:actor_id
func (h *AuditHandler) GetActorActivity(c *gin.Context) {
	actorType := c.Param("actor_type")
	actorID := c.Param("actor_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	} else if limit > 500 {
		limit = 500
	}

	events, err := h.auditService.GetActorActivity(
		c.Request.Context(),
		models.AuditActorType(actorType),
		actorID,
		limit,
	)

	if err != nil {
		log.Error().Err(err).Msg("Error getting actor activity")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo actividad del actor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"actor_type": actorType,
		"actor_id":   actorID,
		"total":      len(events),
		"events":     events,
	})
}

// GetRequestTrace obtiene todos los eventos de un request
// GET /audit/request/:request_id
func (h *AuditHandler) GetRequestTrace(c *gin.Context) {
	requestIDStr := c.Param("request_id")

	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request ID inválido",
		})
		return
	}

	events, err := h.auditService.GetRequestTrace(c.Request.Context(), requestID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting request trace")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo traza del request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id": requestID,
		"total":      len(events),
		"events":     events,
	})
}

// GetRecentEvents obtiene eventos recientes
// GET /audit/recent?hours=24&limit=100
func (h *AuditHandler) GetRecentEvents(c *gin.Context) {
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	if hours <= 0 {
		hours = 24
	} else if hours > 168 { // Máximo 7 días
		hours = 168
	}

	if limit <= 0 {
		limit = 100
	} else if limit > 500 {
		limit = 500
	}

	events, err := h.auditService.GetRecentEvents(c.Request.Context(), hours, limit)
	if err != nil {
		log.Error().Err(err).Msg("Error getting recent events")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo eventos recientes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hours":  hours,
		"total":  len(events),
		"events": events,
	})
}

// Search realiza búsqueda avanzada de eventos
// POST /audit/search
func (h *AuditHandler) Search(c *gin.Context) {
	var req struct {
		EntityType *string    `json:"entity_type"`
		EntityID   *string    `json:"entity_id"`
		Action     *string    `json:"action"`
		ActorType  *string    `json:"actor_type"`
		ActorID    *string    `json:"actor_id"`
		FromDate   *time.Time `json:"from_date"`
		ToDate     *time.Time `json:"to_date"`
		Limit      int        `json:"limit"`
		Offset     int        `json:"offset"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parámetros de búsqueda inválidos",
		})
		return
	}

	// Límites por defecto
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 500 {
		req.Limit = 500
	}

	filter := store.AuditSearchFilter{
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Action:     req.Action,
		ActorType:  req.ActorType,
		ActorID:    req.ActorID,
		FromDate:   req.FromDate,
		ToDate:     req.ToDate,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	events, total, err := h.auditService.Search(c.Request.Context(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Error searching audit events")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error buscando eventos de auditoría",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   total,
		"limit":   req.Limit,
		"offset":  req.Offset,
		"results": len(events),
		"events":  events,
	})
}

// GetStats obtiene estadísticas de auditoría
// GET /audit/stats?hours=24
func (h *AuditHandler) GetStats(c *gin.Context) {
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))
	if hours <= 0 {
		hours = 24
	} else if hours > 168 {
		hours = 168
	}

	events, err := h.auditService.GetRecentEvents(c.Request.Context(), hours, 10000)
	if err != nil {
		log.Error().Err(err).Msg("Error getting audit stats")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo estadísticas",
		})
		return
	}

	// Calcular estadísticas
	stats := map[string]interface{}{
		"total_events": len(events),
		"period_hours": hours,
		"by_action":    make(map[string]int),
		"by_entity":    make(map[string]int),
		"by_actor":     make(map[string]int),
	}

	byAction := make(map[string]int)
	byEntity := make(map[string]int)
	byActor := make(map[string]int)

	for _, event := range events {
		byAction[event.Action]++
		byEntity[event.EntityType]++
		if event.ActorType != nil {
			byActor[*event.ActorType]++
		}
	}

	stats["by_action"] = byAction
	stats["by_entity"] = byEntity
	stats["by_actor"] = byActor

	c.JSON(http.StatusOK, stats)
}

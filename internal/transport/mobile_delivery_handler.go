package transport

import (
	"net/http"

	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MobileDeliveryHandler struct {
	service      service.MobileDeliveryService
	auditService *service.AuditService
}

func NewMobileDeliveryHandler(service service.MobileDeliveryService, auditService *service.AuditService) *MobileDeliveryHandler {
	return &MobileDeliveryHandler{
		service:      service,
		auditService: auditService,
	}
}

// ValidateToken godoc
// @Summary Validar token de entrega del cliente
// @Description El repartidor valida el token proporcionado por el cliente
// @Tags Mobile
// @Accept json
// @Produce json
// @Param request body dto.ValidateTokenRequest true "Token del cliente"
// @Success 200 {object} dto.ValidateTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/mobile/validate-token [post]
func (h *MobileDeliveryHandler) ValidateToken(c *gin.Context) {
	var req dto.ValidateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid validate token request")
		if validationErrors := FormatValidationError(err); len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": err.Error()})
		return
	}

	response, err := h.service.ValidateToken(c.Request.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Error validating token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar el token"})
		return
	}

	// Auditar validación de token
	if h.auditService != nil {
		h.auditService.LogTokenValidated(
			c.Request.Context(),
			req.Token,
			response.Valid,
			c.ClientIP(),
			c.Request.URL.Path,
		)
	}

	if !response.Valid {
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CompleteDelivery godoc
// @Summary Completar entrega
// @Description Marca la entrega como completada y publica mensaje para crear orden de trabajo
// @Tags Mobile
// @Accept json
// @Produce json
// @Param request body dto.MobileCompleteDeliveryRequest true "Información de la entrega completada"
// @Success 200 {object} dto.MobileCompleteDeliveryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/mobile/complete-delivery [post]
func (h *MobileDeliveryHandler) CompleteDelivery(c *gin.Context) {
	var req dto.MobileCompleteDeliveryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid complete delivery request")
		if validationErrors := FormatValidationError(err); len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": err.Error()})
		return
	}

	response, err := h.service.CompleteDelivery(c.Request.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Error completing delivery")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auditar completar entrega
	if h.auditService != nil {
		metadata := map[string]interface{}{
			"validated_dispensers": req.ValidatedDispensers,
			"dispensers_count":     len(req.ValidatedDispensers),
			"work_order_queued":    response.WorkOrderQueued,
		}
		h.auditService.LogDeliveryUpdated(
			c.Request.Context(),
			req.DeliveryID,
			"mobile_app",
			req.Token,
			nil, // before state (opcional)
			response,
			metadata,
		)
	}

	c.JSON(http.StatusOK, response)
}

// SearchDeliveries godoc
// @Summary Buscar deliveries por fecha y reparto
// @Description Búsqueda de deliveries con filtros de fecha (obligatorio) y nro_rto (opcional)
// @Tags Mobile
// @Accept json
// @Produce json
// @Param fecha_accion query string true "Fecha de acción (YYYY-MM-DD)"
// @Param nro_rto query string false "Número de reparto"
// @Success 200 {array} dto.MobileDeliverySearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/mobile/deliveries/search [get]
func (h *MobileDeliveryHandler) SearchDeliveries(c *gin.Context) {
	fechaAccion := c.Query("fecha_accion")
	nroRto := c.Query("nro_rto")

	if fechaAccion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro fecha_accion es obligatorio"})
		return
	}

	results, err := h.service.SearchDeliveries(c.Request.Context(), fechaAccion, nroRto)
	if err != nil {
		log.Error().Err(err).Msg("Error searching deliveries")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

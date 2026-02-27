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
	service service.MobileDeliveryService
}

func NewMobileDeliveryHandler(service service.MobileDeliveryService) *MobileDeliveryHandler {
	return &MobileDeliveryHandler{service: service}
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

	if !response.Valid {
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ValidateDispenser godoc
// @Summary Validar dispenser escaneado
// @Description Valida que el dispenser escaneado pertenezca al delivery
// @Tags Mobile
// @Accept json
// @Produce json
// @Param request body dto.ValidateDispenserRequest true "Información del dispenser"
// @Success 200 {object} dto.ValidateDispenserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/mobile/validate-dispenser [post]
func (h *MobileDeliveryHandler) ValidateDispenser(c *gin.Context) {
	var req dto.ValidateDispenserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid validate dispenser request")
		if validationErrors := FormatValidationError(err); len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": err.Error()})
		return
	}

	response, err := h.service.ValidateDispenser(c.Request.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Error validating dispenser")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar el dispenser"})
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

	c.JSON(http.StatusOK, response)
}

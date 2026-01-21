package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// DeliveryWithTermsHandler maneja endpoints para el flujo integrado de entregas con términos
type DeliveryWithTermsHandler struct {
	service       service.DeliveryWithTermsService
	appBaseURL    string
	termsTTLHours int
}

func NewDeliveryWithTermsHandler(
	service service.DeliveryWithTermsService,
	appBaseURL string,
	termsTTLHours int,
) *DeliveryWithTermsHandler {
	return &DeliveryWithTermsHandler{
		service:       service,
		appBaseURL:    appBaseURL,
		termsTTLHours: termsTTLHours,
	}
}

func (h *DeliveryWithTermsHandler) InitiateDelivery(c *gin.Context) {
	var req dto.InitiateDeliveryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg(constants.LogValidationFailedInitiate)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgInvalidData,
			"message": err.Error(),
		})
		return
	}

	// Validar cantidad de dispensers
	if len(req.Dispensers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgValidationFailed,
			"message": constants.MsgAtLeastOneDispenser,
		})
		return
	}

	if uint(len(req.Dispensers)) != req.Cantidad {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgValidationFailed,
			"message": constants.MsgDispenserQuantityMismatch,
		})
		return
	}

	response, err := h.service.InitiateDelivery(
		c.Request.Context(),
		req,
		h.appBaseURL,
		h.termsTTLHours,
	)
	if err != nil {
		log.Error().Err(err).Msg(constants.LogErrorInitiatingDelivery)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   constants.MsgServerError,
			"message": constants.MsgCouldNotInitiateDelivery,
		})
		return
	}

	log.Info().
		Str("token", response.Token).
		Str("nro_rto", req.NroRto).
		Msg(constants.MsgDeliveryInitiatedSuccess)

	c.JSON(http.StatusOK, response)
}

func (h *DeliveryWithTermsHandler) CompleteDelivery(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgParameterMissing,
			"message": constants.MsgTokenRequired,
		})
		return
	}

	delivery, err := h.service.CompleteDelivery(c.Request.Context(), token)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg(constants.LogErrorCompletingDelivery)

		// Usar helper para determinar código de estado
		statusCode := GetHTTPStatusFromError(err)

		c.JSON(statusCode, gin.H{
			"error":   constants.MsgCouldNotCompleteDelivery,
			"message": err.Error(),
		})
		return
	}

	log.Info().
		Int("delivery_id", delivery.ID).
		Str("nro_rto", delivery.NroRto).
		Str("token", token).
		Msg(constants.MsgDeliveryCompletedSuccess)

	// Convertir a DTO de respuesta usando la función existente
	response := dto.ToDeliveryResponse(delivery)

	c.JSON(http.StatusOK, dto.CompleteDeliveryResponse{
		Success:  true,
		Message:  constants.MsgDeliveryCreatedAfterTerms,
		Delivery: &response,
	})
}

func (h *DeliveryWithTermsHandler) GetDeliveryStatus(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgParameterMissing,
			"message": constants.MsgTokenRequired,
		})
		return
	}

	// Por ahora solo devolvemos un mensaje indicando que deben usar el endpoint de términos
	c.JSON(http.StatusOK, gin.H{
		"message": constants.MsgUseTermsStatusEndpoint,
		"token":   token,
	})
}

package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	service service.DeliveryService
}

func NewDeliveryHandler(service service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: service}
}

func (h *DeliveryHandler) GetAllDeliveries(c *gin.Context) {
	ctx := c.Request.Context()
	nroCta := c.Query("nro_cta")
	fechaStr := c.Query("fecha_accion")
	var deliveries []models.Delivery
	var err error
	if nroCta != "" || fechaStr != "" {
		var fechaAccion *time.Time
		if fechaStr != "" {
			parsed, parseErr := time.Parse("2006-01-02", fechaStr)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha inválida. Formato esperado: YYYY-MM-DD"})
				return
			}
			fechaAccion = &parsed
		}
		deliveries, err = h.service.FindByFilters(ctx, nroCta, fechaAccion)
	} else {
		deliveries, err = h.service.FindAll(ctx)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToDeliveryResponseList(deliveries)
	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) GetDeliveryByID(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}

	delivery, err := h.service.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.MsgDeliveryNotFound})
		return
	}
	response := dto.ToDeliveryResponse(delivery)
	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	ctx := c.Request.Context()
	var delivery models.Delivery

	if err := c.ShouldBindJSON(&delivery); err != nil {
		if validationErrors := FormatValidationError(err); len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": err.Error()})
		return
	}

	if err := h.service.Create(ctx, &delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := dto.ToDeliveryResponse(&delivery)
	c.JSON(http.StatusCreated, response)
}

func (h *DeliveryHandler) UpdateDelivery(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}

	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		if validationErrors := FormatValidationError(err); len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput, "details": err.Error()})
		return
	}

	delivery.ID = id

	if err := h.service.Update(ctx, &delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

func (h *DeliveryHandler) DeleteDelivery(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.MsgDeliveryDeleted})
}

func (h *DeliveryHandler) GetDeliveriesByRto(c *gin.Context) {
	ctx := c.Request.Context()
	nroRto := c.Query("nro_rto")
	fechaStr := c.Query("fecha_accion")

	var fechaAccion *time.Time
	if fechaStr != "" {
		parsed, parseErr := time.Parse("2006-01-02", fechaStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha inválida. Formato esperado: YYYY-MM-DD"})
			return
		}
		fechaAccion = &parsed
	}

	deliveries, err := h.service.FindByRto(ctx, nroRto, fechaAccion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToDeliveryResponseList(deliveries)
	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) GetDeliveriesByNroCta(c *gin.Context) {
	ctx := c.Request.Context()
	nroCta := c.Query("nro_cta")
	fechaStr := c.Query("fecha_accion")

	var fechaAccion *time.Time
	if fechaStr != "" {
		parsed, parseErr := time.Parse("2006-01-02", fechaStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha inválida. Formato esperado: YYYY-MM-DD"})
			return
		}
		fechaAccion = &parsed
	}

	deliveries, err := h.service.FindByFilters(ctx, nroCta, fechaAccion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToDeliveryResponseList(deliveries)
	c.JSON(http.StatusOK, response)
}

// CreateDeliveryFromInfobip maneja la creación de entregas desde el chatbot de Infobip
// POST /api/v1/deliveries/infobip
func (h *DeliveryHandler) CreateDeliveryFromInfobip(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.InfobipDeliveryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := FormatValidationError(err); len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   constants.MsgInvalidInput,
				"details": validationErrors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgInvalidInput,
			"details": err.Error(),
		})
		return
	}

	// Validación adicional de cantidad total
	totalDispensers := req.Tipos.P + req.Tipos.M
	if totalDispensers == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgValidationFailed,
			"message": "Debe especificar al menos un dispenser (P o M)",
		})
		return
	}

	delivery, err := h.service.CreateFromInfobip(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   constants.MsgServerError,
			"message": err.Error(),
		})
		return
	}

	response := dto.InfobipDeliveryResponse{
		Token:   delivery.Token,
		Message: "Entrega creada exitosamente",
	}

	c.JSON(http.StatusCreated, response)
}

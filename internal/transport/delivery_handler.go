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
	service      service.DeliveryService
	auditService *service.AuditService
}

func NewDeliveryHandler(service service.DeliveryService, auditService *service.AuditService) *DeliveryHandler {
	return &DeliveryHandler{
		service:      service,
		auditService: auditService,
	}
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

	// Auditar creación de entrega
	if h.auditService != nil {
		h.auditService.LogDeliveryCreated(
			ctx,
			delivery.ID,
			"api_client",
			"postman",
			&delivery,
			map[string]interface{}{
				"nro_cta":      delivery.NroCta,
				"tipo_entrega": delivery.TipoEntrega,
				"cantidad":     delivery.Cantidad,
			},
		)
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

	// Auditar actualización (sin before state por simplicidad)
	if h.auditService != nil {
		h.auditService.LogDeliveryUpdated(
			ctx,
			delivery.ID,
			"api_client",
			"postman",
			nil, // before state opcional
			&delivery,
			map[string]interface{}{
				"updated_fields": "delivery updated",
			},
		)
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

	// Auditar eliminación
	if h.auditService != nil {
		h.auditService.LogDeliveryUpdated(
			ctx,
			id,
			"api_client",
			"postman",
			nil,
			nil,
			map[string]interface{}{
				"action": "deleted",
			},
		)
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

// GetTallerPrep devuelve los deliveries de una fecha dada con el resumen de dispensers P y M
// que el taller Frio Calor necesita preparar.
// Query param obligatorio: fecha_accion (YYYY-MM-DD)
func (h *DeliveryHandler) GetTallerPrep(c *gin.Context) {
	ctx := c.Request.Context()
	fechaStr := c.Query("fecha_accion")

	if fechaStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'fecha_accion' requerido. Formato: YYYY-MM-DD"})
		return
	}

	if _, err := time.Parse("2006-01-02", fechaStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha inválida. Formato esperado: YYYY-MM-DD"})
		return
	}

	deliveries, err := h.service.FindByFechaAccion(ctx, fechaStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToTallerPrepResponse(fechaStr, deliveries)
	c.JSON(http.StatusOK, response)
}

// GetTokenByFechaAndCta busca y devuelve el token para un delivery específico
// Endpoint público para contact center sin autenticación
// Query params: fecha_accion (YYYY-MM-DD) y nro_cta (string)
func (h *DeliveryHandler) GetTokenByFechaAndCta(c *gin.Context) {
	ctx := c.Request.Context()
	fechaAccion := c.Query("fecha_accion")
	nroCta := c.Query("nro_cta")

	if fechaAccion == "" || nroCta == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros 'fecha_accion' y 'nro_cta' son requeridos"})
		return
	}

	if _, err := time.Parse("2006-01-02", fechaAccion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha inválida. Formato esperado: YYYY-MM-DD"})
		return
	}

	delivery, err := h.service.FindByFechaAndNroCta(ctx, fechaAccion, nroCta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if delivery == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró delivery con los parámetros especificados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           delivery.ID,
		"fecha_accion": delivery.FechaAccion.Format("2006-01-02"),
		"nro_cta":      delivery.NroCta,
		"token":        delivery.Token,
	})
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

package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	service service.DeliveryService
}

func NewDeliveryHandler(service service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: service}
}

func (h *DeliveryHandler) GetAllDeliveries(c *gin.Context) {
	deliveries, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Convertir a DTO y enviar la respuesta a Jmobile o al Chatbot
	response := dto.ToDeliveryResponseList(deliveries)
	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) GetDeliveryByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}

	delivery, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.MsgDeliveryNotFound})
		return
	}
	// Convertir a DTO
	response := dto.ToDeliveryResponse(delivery)
	c.JSON(http.StatusOK, response)
}

func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	var delivery models.Delivery

	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput})
		return
	}

	if err := h.service.Create(&delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convertir a DTO
	response := dto.ToDeliveryResponse(&delivery)
	c.JSON(http.StatusCreated, response)
}

func (h *DeliveryHandler) UpdateDelivery(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}

	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidInput})
		return
	}

	delivery.ID = id

	if err := h.service.Update(&delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

func (h *DeliveryHandler) DeleteDelivery(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.MsgDeliveryDeleted})
}

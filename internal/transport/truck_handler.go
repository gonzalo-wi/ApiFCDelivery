package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TruckHandler struct {
	service service.TruckService
}

func NewTruckHandler(service service.TruckService) *TruckHandler {
	return &TruckHandler{service: service}
}

func (h *TruckHandler) GetAllTrucks(c *gin.Context) {
	ctx := c.Request.Context()
	trucks, err := h.service.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, trucks)
}

func (h *TruckHandler) GetTruckByID(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}
	truck, err := h.service.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, truck)
}

func (h *TruckHandler) CreateTruck(c *gin.Context) {
	ctx := c.Request.Context()
	var truck models.Truck
	if err := c.ShouldBindJSON(&truck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationError(err)})
		return
	}
	err := h.service.Create(ctx, &truck)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Camión creado exitosamente", "data": truck})
}

func (h *TruckHandler) UpdateTruck(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}
	var truck models.Truck
	if err := c.ShouldBindJSON(&truck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationError(err)})
		return
	}
	truck.ID = id
	err = h.service.Update(ctx, &truck)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Camión actualizado exitosamente", "data": truck})
}

func (h *TruckHandler) DeleteTruck(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}
	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.Status(http.StatusNoContent)
}

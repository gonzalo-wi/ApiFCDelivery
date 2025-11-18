package transport

import (
	"GoFrioCalor/internal/constants"
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

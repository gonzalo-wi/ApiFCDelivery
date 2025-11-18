package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DispenserHandler struct {
	service service.DispenserService
}

func NewDispenserHandler(service service.DispenserService) *DispenserHandler {
	return &DispenserHandler{service: service}
}

func (h *DispenserHandler) GetAllDispensers(c *gin.Context) {
	ctx := c.Request.Context()
	dispensers, err := h.service.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dispensers)
}

func (h *DispenserHandler) GetDispenserByID(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}
	dispenser, err := h.service.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, dispenser)
}

func (h *DispenserHandler) CreateDispenser(c *gin.Context) {
	ctx := c.Request.Context()
	var dispenser models.Dispenser
	if err := c.ShouldBindJSON(&dispenser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationError(err)})
		return
	}
	err := h.service.Create(ctx, &dispenser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": constants.MsgDispenserCreated})
}

func (h *DispenserHandler) UpdateDispenser(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}
	var dispenser models.Dispenser
	if err := c.ShouldBindJSON(&dispenser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationError(err)})
		return
	}
	dispenser.ID = id
	err = h.service.Update(ctx, &dispenser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgDispenserUpdated})
}

func (h *DispenserHandler) DeleteDispenser(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidID})
		return
	}
	err = h.service.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgDispenserDeleted})
}

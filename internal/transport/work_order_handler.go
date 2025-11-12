package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkOrderHandler struct {
	pdfService service.PDFService
}

func NewWorkOrderHandler(pdfService service.PDFService) *WorkOrderHandler {
	return &WorkOrderHandler{pdfService: pdfService}
}

func (h *WorkOrderHandler) GenerateWorkOrder(c *gin.Context) {
	var request dto.WorkOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.MsgInvalidData, "details": err.Error()})
		return
	}
	pdfBytes, orderNumber, err := h.pdfService.GenerateWorkOrderPDF(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.MsgPDFGenerationError, "details": err.Error()})
		return
	}
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=orden_trabajo_"+orderNumber+".pdf")
	c.Header("X-Order-Number", orderNumber)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterWorkOrderRoutes(router *gin.RouterGroup, handler *transport.WorkOrderHandler) {
	workOrders := router.Group("/work-orders")
	{
		workOrders.POST("/generate", handler.GenerateWorkOrder)
	}
}

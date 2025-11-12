package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func SetupRouter(deliveryHandler *transport.DeliveryHandler, workOrderHandler *transport.WorkOrderHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		RegisterDeliveryRoutes(api, deliveryHandler)
		RegisterWorkOrderRoutes(api, workOrderHandler)
	}

	return router
}

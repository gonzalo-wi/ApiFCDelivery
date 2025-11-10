package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterDeliveryRoutes(router *gin.RouterGroup, handler *transport.DeliveryHandler) {
	deliveries := router.Group("/deliveries")
	{
		deliveries.GET("", handler.GetAllDeliveries)
		deliveries.GET("/:id", handler.GetDeliveryByID)
		deliveries.POST("", handler.CreateDelivery)
		deliveries.PUT("/:id", handler.UpdateDelivery)
		deliveries.DELETE("/:id", handler.DeleteDelivery)
	}
}

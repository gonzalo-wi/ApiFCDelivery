package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterMobileRoutes(rg *gin.RouterGroup, handler *transport.MobileDeliveryHandler) {
	mobile := rg.Group("/mobile")
	{
		mobile.POST("/validate-token", handler.ValidateToken)

		mobile.POST("/complete-delivery", handler.CompleteDelivery)

		mobile.GET("/deliveries/search", handler.SearchDeliveries)
	}
}

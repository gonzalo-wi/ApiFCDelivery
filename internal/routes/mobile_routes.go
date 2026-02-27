package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterMobileRoutes(rg *gin.RouterGroup, handler *transport.MobileDeliveryHandler) {
	mobile := rg.Group("/mobile")
	{
		mobile.POST("/validate-token", handler.ValidateToken)
		mobile.POST("/validate-dispenser", handler.ValidateDispenser)
		mobile.POST("/complete-delivery", handler.CompleteDelivery)
	}
}

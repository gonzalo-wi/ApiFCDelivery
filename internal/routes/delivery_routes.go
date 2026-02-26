package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterDeliveryRoutes(router *gin.RouterGroup, handler *transport.DeliveryHandler) {
	deliveries := router.Group("/deliveries")
	{
		deliveries.GET("", handler.GetAllDeliveries)
		deliveries.GET("/by-rto", handler.GetDeliveriesByRto)
		deliveries.GET("/by-cta", handler.GetDeliveriesByNroCta)
		deliveries.GET("/:id", handler.GetDeliveryByID)
		deliveries.POST("", handler.CreateDelivery)
		deliveries.POST("/infobip", handler.CreateDeliveryFromInfobip)
		deliveries.PUT("/:id", handler.UpdateDelivery)
		deliveries.DELETE("/:id", handler.DeleteDelivery)
	}
}

// RegisterDeliveryWithTermsRoutes registra rutas para el flujo integrado con términos y condiciones
func RegisterDeliveryWithTermsRoutes(router *gin.RouterGroup, handler *transport.DeliveryWithTermsHandler) {
	deliveries := router.Group("/deliveries")
	{
		// Flujo nuevo: Iniciar entrega → Aceptar términos → Completar entrega
		deliveries.POST("/initiate", handler.InitiateDelivery)
		deliveries.POST("/complete/:token", handler.CompleteDelivery)
		deliveries.GET("/status/:token", handler.GetDeliveryStatus)
	}
}

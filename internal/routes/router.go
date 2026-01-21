package routes

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/middleware"
	"GoFrioCalor/internal/transport"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(deliveryHandler *transport.DeliveryHandler, dispenserHandler *transport.DispenserHandler,
	workOrderHandler *transport.WorkOrderHandler, termsSessionHandler *transport.TermsSessionHandler,
	deliveryWithTermsHandler *transport.DeliveryWithTermsHandler, cfg *config.Config) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(middleware.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.GetCORSOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		RegisterDeliveryRoutes(api, deliveryHandler)
		RegisterDispenserRoutes(api, dispenserHandler)
		RegisterWorkOrderRoutes(api, workOrderHandler)
		RegisterTermsRoutes(api, termsSessionHandler)
		RegisterDeliveryWithTermsRoutes(api, deliveryWithTermsHandler)
	}
	return router
}

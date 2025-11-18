package routes

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/transport"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(deliveryHandler *transport.DeliveryHandler, dispenserHandler *transport.DispenserHandler,
	workOrderHandler *transport.WorkOrderHandler, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Configuración de CORS
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
	}

	return router
}

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

	// Deshabilitar el redirect automático de trailing slashes
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

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

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

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

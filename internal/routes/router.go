package routes

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/middleware"
	"GoFrioCalor/internal/transport"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(deliveryHandler *transport.DeliveryHandler,
	workOrderHandler *transport.WorkOrderHandler, termsSessionHandler *transport.TermsSessionHandler,
	deliveryWithTermsHandler *transport.DeliveryWithTermsHandler, mobileDeliveryHandler *transport.MobileDeliveryHandler,
	auditHandler *transport.AuditHandler, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Deshabilitar el redirect automático de trailing slashes
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(gin.Recovery())

	// Request ID middleware (para trazabilidad)
	router.Use(middleware.RequestID())

	router.Use(middleware.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.GetCORSOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-api-key", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
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

	// Rutas públicas de autenticación (sin middleware)
	authHandler := transport.NewAuthHandler(cfg.AuthServiceURL)
	RegisterAuthRoutes(router, authHandler)

	// Grupo de rutas protegidas con autenticación JWT
	api := router.Group("/dispenser-operations/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.AuthServiceURL))
	{
		RegisterDeliveryRoutes(api, deliveryHandler)
		RegisterWorkOrderRoutes(api, workOrderHandler)
		RegisterTermsRoutes(api, termsSessionHandler)
		RegisterDeliveryWithTermsRoutes(api, deliveryWithTermsHandler)

		// Mobile routes (solo si el handler está disponible)
		if mobileDeliveryHandler != nil {
			RegisterMobileRoutes(api, mobileDeliveryHandler)
		}

		// Audit routes (si el handler está disponible)
		if auditHandler != nil {
			RegisterAuditRoutes(api, auditHandler)
		}
	}
	return router
}

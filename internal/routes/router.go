package routes

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/middleware"
	"GoFrioCalor/internal/transport"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(deliveryHandler *transport.DeliveryHandler,
	workOrderHandler *transport.WorkOrderHandler, termsSessionHandler *transport.TermsSessionHandler,
	deliveryWithTermsHandler *transport.DeliveryWithTermsHandler, mobileDeliveryHandler *transport.MobileDeliveryHandler,
	auditHandler *transport.AuditHandler, cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.PrometheusMetrics())
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

	// Debug endpoint para listar todas las rutas registradas
	router.GET("/debug/routes", func(c *gin.Context) {
		routes := router.Routes()
		c.JSON(200, routes)
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	authHandler := transport.NewAuthHandler(cfg.AuthServiceURL)
	RegisterAuthRoutes(router, authHandler)

	// ===== RUTAS PÚBLICAS (SIN AUTENTICACIÓN) =====
	// IMPORTANTE: Se registran directamente en el router para que tengan máxima prioridad
	// y no sean capturadas por rutas con parámetros dinámicos como /:id
	router.GET("/dispenser-operations/api/v1/deliveries/taller-prep", deliveryHandler.GetTallerPrep)
	router.GET("/dispenser-operations/api/v1/deliveries/contact-center/token", deliveryHandler.GetTokenByFechaAndCta)
	router.POST("/dispenser-operations/api/v1/deliveries/contact-center", deliveryHandler.CreateDeliveryFromContactCenter)

	// Rutas públicas de términos (para que el cliente pueda aceptar/rechazar)
	publicAPI := router.Group("/dispenser-operations/api/v1")
	RegisterPublicTermsRoutes(publicAPI, termsSessionHandler)

	// ===== RUTAS PROTEGIDAS (CON AUTENTICACIÓN) =====
	api := router.Group("/dispenser-operations/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.AuthServiceURL))
	{
		RegisterDeliveryRoutes(api, deliveryHandler)
		RegisterWorkOrderRoutes(api, workOrderHandler)
		RegisterTermsRoutes(api, termsSessionHandler)
		RegisterDeliveryWithTermsRoutes(api, deliveryWithTermsHandler)

		if mobileDeliveryHandler != nil {
			RegisterMobileRoutes(api, mobileDeliveryHandler)
		}

		if auditHandler != nil {
			RegisterAuditRoutes(api, auditHandler)
		}
	}
	return router
}

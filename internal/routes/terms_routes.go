package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

// RegisterTermsRoutes registra rutas protegidas de términos (requieren autenticación)
func RegisterTermsRoutes(router *gin.RouterGroup, handler *transport.TermsSessionHandler) {
	infobip := router.Group("/infobip")
	{
		infobip.POST("/session", handler.CreateInfobipSession)
	}
}

// RegisterPublicTermsRoutes registra rutas públicas de términos (sin autenticación)
// Incluye: creación de sesión desde contact center, consulta de estado y aceptación/rechazo por cliente
func RegisterPublicTermsRoutes(router *gin.RouterGroup, handler *transport.TermsSessionHandler) {
	// Contact center - crear sesión sin autenticación (uso interno)
	contactCenter := router.Group("/contact-center")
	{
		contactCenter.POST("/session", handler.CreateContactCenterSession)
	}

	terms := router.Group("/terms")
	{
		// Cliente final - consultar y aceptar/rechazar términos
		terms.GET("/:token", handler.GetTermsStatus)
		terms.POST("/:token/accept", handler.AcceptTerms)
		terms.POST("/:token/reject", handler.RejectTerms)

		// Panel contact center - consultar estado por sessionId
		terms.GET("/by-session/:sessionId", handler.GetSessionBySessionID)
	}
}

package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterTermsRoutes(router *gin.RouterGroup, handler *transport.TermsSessionHandler) {
	infobip := router.Group("/infobip")
	{
		infobip.POST("/session", handler.CreateInfobipSession)
	}
	terms := router.Group("/terms")
	{
		terms.GET("/:token", handler.GetTermsStatus)
		terms.GET("/by-session/:sessionId", handler.GetSessionBySessionID)
		terms.POST("/:token/accept", handler.AcceptTerms)
		terms.POST("/:token/reject", handler.RejectTerms)
	}
}

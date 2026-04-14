package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, authHandler *transport.AuthHandler) {
	auth := router.Group("/dispenser-operations/auth")
	{
		auth.GET("/generar-token", authHandler.GenerateToken)
	}
}

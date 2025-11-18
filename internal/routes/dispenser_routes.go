package routes

import (
	"GoFrioCalor/internal/transport"

	"github.com/gin-gonic/gin"
)

func RegisterDispenserRoutes(router *gin.RouterGroup, handler *transport.DispenserHandler) {
	dispensers := router.Group("/dispensers")
	{
		dispensers.GET("", handler.GetAllDispensers)
		dispensers.GET("/:id", handler.GetDispenserByID)
		dispensers.POST("", handler.CreateDispenser)
		dispensers.PUT("/:id", handler.UpdateDispenser)
		dispensers.DELETE("/:id", handler.DeleteDispenser)
	}
}

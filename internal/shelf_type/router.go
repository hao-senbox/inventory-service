package shelftype

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ShelfTypeHandler) {
	api := r.Group("api/v1")
	{
		location := api.Group("/shelf_type").Use(middleware.Secured())
		{
			location.POST("", handler.CreateShelfType)
			location.GET("", handler.GetShelfTypes)
			location.GET("/:id", handler.GetShelfTypeByID)
			location.PUT("/:id", handler.UpdateShelfType)
			location.DELETE("/:id", handler.DeleteShelfType)
		}
	}
}

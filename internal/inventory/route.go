package inventory

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *InventoryHandler) {
	
	api := r.Group("api/v1") 
	{
		location := api.Group("/locations").Use(middleware.Secured())
		{
			location.POST("", handler.CreateLocation)
			location.GET("", handler.GetLocations)
			location.GET("/:id", handler.GetLocationByID)
		}
		// inventory := api.Group("/inventory").Use(middleware.Secured())
		// {
		// 	inventory.POST("", handler.AddInventory)
		// }
	}
}
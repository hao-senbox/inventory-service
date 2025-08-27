package location

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *LocationHandler) {
	api := r.Group("api/v1")
	{
		location := api.Group("/location").Use(middleware.Secured())
		{
			location.POST("", handler.CreateLocation)
			location.GET("", handler.GetLocations)
			location.GET("/:id", handler.GetLocationByID)
			location.PUT("/:id", handler.UpdateLocation)
			location.DELETE("/:id", handler.DeleteLocation)
			location.PUT("add-product/:id", handler.AddProductToLocation)
		}
	}
}

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
			// location.GET("", handler.GetLocations)
			// location.GET("/:id", handler.GetLocationByID)
		}
	}
}
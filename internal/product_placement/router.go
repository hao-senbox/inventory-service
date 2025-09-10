package productplacement

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ProductPlacementHandler) {
	api := r.Group("api/v1")
	{
		location := api.Group("/product_placement").Use(middleware.Secured())
		{
			location.GET("shelf/:id", handler.GetProductPlacementsByShelfID)
			location.GET("product/:id", handler.GetProductPlacementsByProductID)
		}
	}
}
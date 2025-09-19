package shelfquantity

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ShelfQuantityHandler) {
	api := r.Group("api/v1")
	{
		location := api.Group("/shelf_quantity").Use(middleware.Secured())
		{
			location.POST("", handler.CreateShelfQuantity)
			location.GET("/shelf/:id", handler.GetShelfQuantitiesByShelfID)
			// location.GET("/product/:id", handler.GetShelfQuantitiesByProductID)
		}
	}
}

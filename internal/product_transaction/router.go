package producttransaction

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ProductTransactionHandler) {
	api := r.Group("api/v1")
	{
		location := api.Group("/product_transaction").Use(middleware.Secured())
		{
			location.POST("", handler.CreateProductTransaction)
			// location.GET("", handler.GetProductTransactions)
			// location.GET("/:id", handler.GetProductTransactionByID)
			// location.PUT("/:id", handler.UpdateProductTransaction)
			// location.DELETE("/:id", handler.DeleteProductTransaction)
		}
	}
}
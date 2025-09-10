package storage

import (
	"inventory-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *StorageHandler) {

	api := r.Group("api/v1")
	{
		location := api.Group("/storage").Use(middleware.Secured())
		{
			location.POST("", handler.CreateStorage)
			location.GET("", handler.GetStoragies)
			location.GET("/:id", handler.GetStorageByID)
			location.PUT("/:id", handler.UpdateStorage)
			location.DELETE("/:id", handler.DeleteStorage)
			location.GET("/tree", handler.GetStorageTree)
		}
	}
}

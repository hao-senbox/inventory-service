package storage

import (
	"context"
	"fmt"
	"inventory-service/helper"
	"inventory-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	StorageService StorageService
}

func NewStorageHandler(storageService StorageService) *StorageHandler {
	return &StorageHandler{
		StorageService: storageService,
	}
}

func (h *StorageHandler) CreateStorage(c *gin.Context) {
	
	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	var req CreateStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	storageID, err := h.StorageService.CreateStorage(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", storageID)

}

func (h *StorageHandler) GetStoragies(c *gin.Context) {

	typeString := c.Query("type")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	storagies, err := h.StorageService.GetStoragies(ctx, typeString)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", storagies)

}

func (h *StorageHandler) GetStorageByID(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	storage, err := h.StorageService.GetStorageByID(ctx, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", storage)
}

func (h *StorageHandler) GetStorageTree(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	storageTree, err := h.StorageService.GetStorageTree(ctx)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", storageTree)
	
}
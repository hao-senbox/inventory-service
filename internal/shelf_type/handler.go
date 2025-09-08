package shelftype

import (
	"context"
	"fmt"
	"inventory-service/helper"
	"inventory-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type ShelfTypeHandler struct {
	ShelfTypeService ShelfTypeService
}

func NewShelfTypeHandler(shelfTypeService ShelfTypeService) *ShelfTypeHandler {
	return &ShelfTypeHandler{
		ShelfTypeService: shelfTypeService,
	}
}

func (h *ShelfTypeHandler) CreateShelfType(c *gin.Context) {

	var req CreateShelfTypeRequest
	
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

	shelfType, err := h.ShelfTypeService.CreateShelfType(ctx, &req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Add shelf type successfully", shelfType)

}

func (h *ShelfTypeHandler) GetShelfTypes(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	shelfTypes, err := h.ShelfTypeService.GetShelfTypes(ctx)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get shelf types successfully", shelfTypes)

}

func (h *ShelfTypeHandler) GetShelfTypeByID(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	shelfType, err := h.ShelfTypeService.GetShelfTypeByID(ctx, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get shelf type successfully", shelfType)

}
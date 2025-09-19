package shelfquantity

import (
	"context"
	"fmt"
	"inventory-service/helper"
	"inventory-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type ShelfQuantityHandler struct {
	ShelfQuantityService ShelfQuantityService
}

func NewShelfQuantityHandler(service ShelfQuantityService) *ShelfQuantityHandler {
	return &ShelfQuantityHandler{
		ShelfQuantityService: service,
	}
}

func (h *ShelfQuantityHandler) CreateShelfQuantity(c *gin.Context) {

	var req []CreateShelfQuantityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ShelfQuantityService.CreateShelfQuantity(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create shelf quantity successfully", nil)

}

func (h *ShelfQuantityHandler) GetShelfQuantitiesByShelfID(c *gin.Context) {

	shelfID := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	shelfQuantities, err := h.ShelfQuantityService.GetShelfQuantitiesByShelfID(ctx, shelfID)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get shelf quantity successfully", shelfQuantities)
	
}
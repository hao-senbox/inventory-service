package producttransaction

import (
	"context"
	"fmt"
	"inventory-service/helper"
	"inventory-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type ProductTransactionHandler struct {
	ProductTransactionService ProductTransactionService
}

func NewProductTransactionHandler(productTransactionService ProductTransactionService) *ProductTransactionHandler {
	return &ProductTransactionHandler{
		ProductTransactionService: productTransactionService,
	}
}

func (h *ProductTransactionHandler) CreateProductTransaction(c *gin.Context) {

	var req CreateProductTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID , exists := c.Get(constants.UserID)
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

	productTransactionID, err := h.ProductTransactionService.CreateProductTransaction(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create product transaction successfully", productTransactionID)
	
}

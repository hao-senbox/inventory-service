package productplacement

import (
	"inventory-service/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductPlacementHandler struct {
	ProductPlacementService ProductPlacementService
}

func NewProductPlacementHandler(productPlacementService ProductPlacementService) *ProductPlacementHandler {
	return &ProductPlacementHandler{
		ProductPlacementService: productPlacementService,
	}
}

func (h *ProductPlacementHandler) GetProductPlacementsByShelfID(c *gin.Context) {

	shelfId := c.Param("shelf_id")

	placements, err := h.ProductPlacementService.GetProductPlacementsByShelfID(c, shelfId)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get product placements by shelf id successfully", placements)

}

func (h *ProductPlacementHandler) GetProductPlacementsByProductID(c *gin.Context) {

	productID := c.Param("product_id")

	placements, err := h.ProductPlacementService.GetProductPlacementsByProductID(c, productID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get product placements by product id successfully", placements)

}

package location

import (
	"context"
	"fmt"
	"inventory-service/helper"
	"inventory-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	service LocationService
}

func NewLocationHandler(service LocationService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {

	var req CreateLocationRequest
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

	locationID, err := h.service.CreateLocation(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Create location successfully", locationID)

}

func (h *LocationHandler) GetLocations(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	locations, err := h.service.GetLocations(ctx)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get location successfully", locations)

}

func (h *LocationHandler) GetLocationByID(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	location, err := h.service.GetLocationByID(ctx, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get location successfully", location)

}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {

	id := c.Param("id")

	var req UpdateLocationRequest
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

	err := h.service.UpdateLocation(ctx, &req, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Update location successfully", nil)

}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.service.DeleteLocation(ctx, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Delete location successfully", nil)
	
}

func (h *LocationHandler) AddProductToLocation(c *gin.Context) {

	id := c.Param("id")

	var req AddProductToLocationRequest
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

	err := h.service.AddProductToLocation(ctx, &req, id)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Add product to location successfully", nil)

}
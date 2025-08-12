package inventory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateLocationRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	ImageMain   string  `json:"image_main"`
	ImageMap    string  `json:"image_map"`
	ParentID    *string `json:"parent_id"`
}

type AddInventoryRequest struct {
	ProductID   string             `json:"product_id" binding:"required"`
	LocationID  primitive.ObjectID `json:"location_id" binding:"required"`
	Quantity    int                `json:"quantity" binding:"required,min=1"`
	BatchNumber string             `json:"batch_number"`
	ExpiryDate  *time.Time         `json:"expiry_date"`
}

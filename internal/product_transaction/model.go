package producttransaction

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductTransaction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	ShelfID   primitive.ObjectID `json:"shelf_id" bson:"shelf_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Action    string             `json:"action" bson:"action"`
	ActionBy  string             `json:"action_by" bson:"action_by"`
	ActionAt  time.Time          `json:"action_at" bson:"action_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateProductPlacementRequest struct {
	ProductID   string `json:"product_id"`
	ShelfID     string `json:"shelf_id"`
	CurrentQty  int    `json:"current_qty"`
}
package productplacement

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductPlacement struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	ProductID   primitive.ObjectID   `json:"product_id" bson:"product_id"`
	ShelfID     primitive.ObjectID   `json:"shelf_id" bson:"shelf_id"`
	CurrentQty  int                  `json:"current_qty" bson:"current_qty"`
	Path        string               `json:"path" bson:"path"`
	AncestorIDs []primitive.ObjectID `json:"ancestor_ids" bson:"ancestor_ids"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

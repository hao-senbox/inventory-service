package shelfquantity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShelfQuantity struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ShelfID   primitive.ObjectID `json:"shelf_id" bson:"shelf_id"`
	Code      string             `json:"code" bson:"code"`
	QRCode    string             `json:"qrcode" bson:"qrcode"`
	Note      string             `json:"note" bson:"note"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

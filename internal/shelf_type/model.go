package shelftype

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShelfType struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Note      *string            `json:"note" bson:"note"`
	Slot      *int               `json:"slot" bson:"slot"`
	Level     *int               `json:"level" bson:"level"`
	Stock     *int               `json:"stock" bson:"stock"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

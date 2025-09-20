package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	QRCode      string               `json:"qrcode" bson:"qrcode"`
	Type        string               `json:"type" bson:"type"` // "warehouse", "building", "floor", "room", "shelf"
	Description *string              `json:"description" bson:"description"`
	ImageMain   *string              `json:"main_image" bson:"main_image"`
	ImageMap    *string              `json:"map_image" bson:"map_image"`
	ParentID    *primitive.ObjectID  `json:"parent_id" bson:"parent_id"`
	AncestorIDs []primitive.ObjectID `json:"ancestor_ids" bson:"ancestor_ids"`
	Level       int                  `json:"level" bson:"level"`
	Path        string               `json:"path" bson:"path"`
	IsActive    bool                 `json:"is_actice" bson:"is_actice"`

	ShelfTypeID *primitive.ObjectID `json:"shelf_type_id,omitempty" bson:"shelf_type_id,omitempty"`
	ShelfID     *string             `json:"shelf_id" bson:"shelf_id"`
	Slots       *int                `json:"slots,omitempty" bson:"slots,omitempty"`
	Levels      *int                `json:"levels,omitempty" bson:"levels,omitempty"`
	TotalStock  *int                `json:"total_stock,omitempty" bson:"total_stock,omitempty"`

	CreatedBy string    `json:"created_by" bson:"created_by"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

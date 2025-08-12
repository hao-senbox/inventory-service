package inventory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Locations struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	QRCode      string               `json:"qrcode" bson:"qrcode"`
	Type        string               `json:"type" bson:"type"`
	Description *string              `json:"description" bson:"description"`
	ImageMain   string               `json:"image_main" bson:"image"`
	ImageMap    string               `json:"image_map" bson:"image_map"`
	ParentID    *primitive.ObjectID  `json:"parent_id" bson:"parent_id"`
	AncestorIDs []primitive.ObjectID `json:"ancestor_ids" bson:"ancestor_ids"`
	Level       int                  `json:"level" bson:"level"`
	Path        string               `json:"path" bson:"path"`
	IsActice    bool                 `json:"is_actice" bson:"is_actice"`
	CreatedBy   string               `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

type LocationHierarchy struct {
	WarehouseID      *primitive.ObjectID `bson:"warehouse_id,omitempty" json:"warehouse_id,omitempty"`
	WarehouseName    string              `bson:"warehouse_name" json:"warehouse_name"`
	BuildingID       *primitive.ObjectID `bson:"building_id,omitempty" json:"building_id,omitempty"`
	BuildingName     string              `bson:"building_name" json:"building_name"`
	FloorID          *primitive.ObjectID `bson:"floor_id,omitempty" json:"floor_id,omitempty"`
	FloorName        string              `bson:"floor_name" json:"floor_name"`
	RoomID           *primitive.ObjectID `bson:"room_id,omitempty" json:"room_id,omitempty"`
	RoomName         string              `bson:"room_name" json:"room_name"`
	ShelfID          *primitive.ObjectID `bson:"shelf_id,omitempty" json:"shelf_id,omitempty"`
	ShelfName        string              `bson:"shelf_name" json:"shelf_name"`
	ShelfFloorID     *primitive.ObjectID `bson:"shelf_floor_id,omitempty" json:"shelf_floor_id,omitempty"`
	ShelfFloorName   string              `bson:"shelf_floor_name" json:"shelf_floor_name"`
	ShelfSlotID      *primitive.ObjectID `bson:"shelf_slot_id,omitempty" json:"shelf_slot_id,omitempty"`
	ShelfSlotName    string              `bson:"shelf_slot_name" json:"shelf_slot_name"`
	ShelfSlotBoxID   *primitive.ObjectID `bson:"shelf_slot_box_id,omitempty" json:"shelf_slot_box_id,omitempty"`
	ShelfSlotBoxName string              `bson:"shelf_slot_box_name" json:"shelf_slot_box_name"`
}

type InventoryItem struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID         string             `bson:"product_id" json:"product_id"`
	LocationID        primitive.ObjectID `bson:"location_id" json:"location_id"`
	LocationHierarchy LocationHierarchy  `bson:"location_hierarchy" json:"location_hierarchy"`
	Quantity          int                `bson:"quantity" json:"quantity"`
	BatchNumber       string             `bson:"batch_number" json:"batch_number"`
	ExpiryDate        *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CreatedBy         primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	IsActive          bool               `bson:"is_active" json:"is_active"`
}

type InventoryHistory struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ActionType       string              `bson:"action_type" json:"action_type"` // ADD, MOVE, REMOVE, UPDATE
	ProductID        string              `bson:"product_id" json:"product_id"`
	Quantity         int                 `bson:"quantity" json:"quantity"`
	FromLocationID   *primitive.ObjectID `bson:"from_location_id,omitempty" json:"from_location_id,omitempty"`
	ToLocationID     *primitive.ObjectID `bson:"to_location_id,omitempty" json:"to_location_id,omitempty"`
	FromLocationPath string              `bson:"from_location_path" json:"from_location_path"`
	ToLocationPath   string              `bson:"to_location_path" json:"to_location_path"`
	PerformedBy      primitive.ObjectID  `bson:"performed_by" json:"performed_by"`
	PerformedAt      time.Time           `bson:"performed_at" json:"performed_at"`
	Notes            string              `bson:"notes" json:"notes"`
	Metadata         map[string]any      `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

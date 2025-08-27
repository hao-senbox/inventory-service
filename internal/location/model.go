package location

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id"`
	LocationName    string               `json:"location_name" bson:"location_name"`
	LocationAddress string               `json:"location_address" bson:"location_address"`
	LocationImage   string               `json:"location_image" bson:"location_image"`
	LocationMap     string               `json:"location_map" bson:"location_map"`
	LocationQRCode  string               `json:"location_qrcode" bson:"location_qrcode"`
	WareHouseID     primitive.ObjectID   `json:"warehouse_id" bson:"warehouse_id"`
	BuildingID      primitive.ObjectID   `json:"building_id" bson:"building_id"`
	FloorID         primitive.ObjectID   `json:"floor_id" bson:"floor_id"`
	RoomID          primitive.ObjectID   `json:"room_id" bson:"room_id"`
	ShelfID         []primitive.ObjectID `json:"shelf_id" bson:"shelf_id"`
	Products        []Product            `json:"products" bson:"products"`
	IsActice        bool                 `json:"is_actice" bson:"is_actice"`
	CreatedBy       string               `json:"created_by" bson:"created_by"`
	CreatedAt       time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at" bson:"updated_at"`
}

type Product struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}	
package location

type CreateLocationRequest struct {
	LocationName    string   `json:"location_name" bson:"location_name"`
	LocationAddress string   `json:"location_address" bson:"location_address"`
	LocationImage   string   `json:"location_image" bson:"location_image"`
	LocationMap     string   `json:"location_map" bson:"location_map"`
	WareHouseID     string   `json:"warehouse_id" bson:"warehouse_id"`
	BuildingID      string   `json:"building_id" bson:"building_id"`
	FloorID         string   `json:"floor_id" bson:"floor_id"`
	RoomID          string   `json:"room_id" bson:"room_id"`
	ShelfID         []string `json:"shelf_id" bson:"shelf_id"`
	CreatedBy       string   `json:"created_by" bson:"created_by"`
}

type UpdateLocationRequest struct {
	LocationName    string   `json:"location_name" bson:"location_name"`
	LocationAddress string   `json:"location_address" bson:"location_address"`
	LocationImage   string   `json:"location_image" bson:"location_image"`
	LocationMap     string   `json:"location_map" bson:"location_map"`
	WareHouseID     string   `json:"warehouse_id" bson:"warehouse_id"`
	BuildingID      string   `json:"building_id" bson:"building_id"`
	FloorID         string   `json:"floor_id" bson:"floor_id"`
	RoomID          string   `json:"room_id" bson:"room_id"`
	ShelfID         []string `json:"shelf_id" bson:"shelf_id"`
	UpdatedBy       string   `json:"updated_by" bson:"updated_by"`
}

type AddProductToLocationRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}
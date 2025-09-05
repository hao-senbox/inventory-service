package location

import "time"

type LocationResponse struct {
	ID              string         `json:"location_id"`
	LocationName    string         `json:"location_name"`
	LocationAddress string         `json:"location_address"`
	QrCode          string         `json:"location_qrcode"`
	LocationImage   string         `json:"location_image"`
	LocationMap     string         `json:"location_map"`
	Warehouse       *WarehouseInfo `json:"warehouse"`
	Building        *BuildingInfo  `json:"building"`
	Floor           *FloorInfo     `json:"floor"`
	Room            *RoomInfo      `json:"room"`
	Shelves         []*ShelfInfo   `json:"shelves"`
	Products        []*ProductInfo `json:"products"`
	CreatedBy       string         `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	IsActive        bool           `json:"is_active"`
}

type WarehouseInfo struct {
	ID   string `json:"warehouse_id"`
	Name string `json:"warehouse_name"`
}

type BuildingInfo struct {
	ID   string `json:"building_id"`
	Name string `json:"building_name"`
}

type FloorInfo struct {
	ID   string `json:"floor_id"`
	Name string `json:"floor_name"`
}

type RoomInfo struct {
	ID   string `json:"room_id"`
	Name string `json:"room_name"`
}

type ShelfInfo struct {
	ID   string `json:"shelf_id"`
	Name string `json:"shelf_name"`
}

type ProductInfo struct {
	ID           string  `json:"id" bson:"_id"`
	Name         string  `json:"name" bson:"name"`
	PriceStore   float64 `json:"price_store" bson:"price_store"`
	PriceService float64 `json:"price_service" bson:"price_service"`
	Description  string  `json:"description" bson:"description"`
	Image        string  `json:"image" bson:"image"`
	FolderName   string  `json:"folder_name" bson:"folder_name"`
	TopicName    string  `json:"topic_name" bson:"topic_name"`
	Quantity     int     `json:"quantity"`
}

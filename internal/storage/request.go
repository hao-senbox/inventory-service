package storage

type CreateStorageRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	ImageMain   *string  `json:"main_image"`
	ImageMap    *string  `json:"map_image"`
	ParentID    *string `json:"parent_id"`
	ShelfTypeID *string `json:"shelf_type_id"`
	Slots       *int    `json:"slots"`
	Levels      *int    `json:"levels"`
}

type UpdateStorageRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string  `json:"description"`
	ImageMain   *string  `json:"main_image"`
	ImageMap    *string  `json:"map_image"`
	ParentID    *string `json:"parent_id"`
	ShelfTypeID *string `json:"shelf_type_id"`
	Slots       *int    `json:"slots"`
	Levels      *int    `json:"levels"`
}

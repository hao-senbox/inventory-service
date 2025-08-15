package storage


type CreateStorageRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	ImageMain   string  `json:"image_main"`
	ImageMap    string  `json:"image_map"`
	ParentID    *string `json:"parent_id"`
}

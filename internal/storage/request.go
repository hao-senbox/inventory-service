package storage


type CreateStorageRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	ImageMain   string  `json:"main_image"`
	ImageMap    string  `json:"map_image"`
	ParentID    *string `json:"parent_id"`
}

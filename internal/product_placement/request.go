package productplacement

type CreateProductPlacementRequest struct {
	ProductID   string `json:"product_id"`
	ShelfID     string `json:"shelf_id"`
	CurrentQty  int    `json:"current_qty"`
}

type UpdateProductPlacementRequest struct {
	ProductID   string `json:"product_id"`
	ShelfID     string `json:"shelf_id"`
	CurrentQty  int    `json:"current_qty"`
}
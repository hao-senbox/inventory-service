package shelfquantity

type CreateShelfQuantityRequest struct {
	ShelfID string `json:"shelf_id" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Note    string `json:"note"`
}

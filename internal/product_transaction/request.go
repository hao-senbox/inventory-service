package producttransaction

type CreateProductTransactionRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	ShelfID   string `json:"shelf_id" bson:"shelf_id"`
	Quantity  int    `json:"quantity" bson:"quantity"`
	Action    string `json:"action" bson:"action"`
}

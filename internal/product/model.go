package product

type Product struct {
	ID           string  `json:"id" bson:"_id"`
	Name         string  `json:"name" bson:"name"`
	PriceStore   float64 `json:"price_store" bson:"price_store"`
	PriceService float64 `json:"price_service" bson:"price_service"`
	Description  string  `json:"description" bson:"description"`
	Image        string  `json:"image" bson:"image"`
	FolderName   string  `json:"folder_name" bson:"folder_name"`
	TopicName    string  `json:"topic_name" bson:"topic_name"`
}

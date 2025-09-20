package ports

import (
	"context"
	"inventory-service/internal/shared/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetStorageByShelfID(ctx context.Context, id primitive.ObjectID) ([]*model.Storage, error)
	DeleteStorage(ctx context.Context, id primitive.ObjectID) error
	GetStorageByID(ctx context.Context, id *primitive.ObjectID) (*model.Storage, error)
	CheckShelfType(ctx context.Context, shelf_type_id primitive.ObjectID) (bool, error)
}

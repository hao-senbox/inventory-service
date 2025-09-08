package ports

import (
	"context"
	"inventory-service/internal/shared/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetStorageByShelfID(ctx context.Context, id primitive.ObjectID) ([]*model.Storage, error)
	DeleteStorage(ctx context.Context, id primitive.ObjectID) error
}

package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryRepository interface {
	AddStorage(ctx context.Context, storage *Storage) (string, error)
	GetStoragies(ctx context.Context, typeString string) ([]*Storage, error)
	GetStorageByID(ctx context.Context, id *primitive.ObjectID) (*Storage, error)
}

type inventoryRepository struct {
	storageCollection *mongo.Collection
}

func NewInventoryRepository(storageCollection *mongo.Collection) InventoryRepository {
	return &inventoryRepository{
		storageCollection: storageCollection,
	}
}

func (r *inventoryRepository) AddStorage(ctx context.Context, storage *Storage) (string, error) {

	result, err := r.storageCollection.InsertOne(ctx, storage)
	if err != nil {
		return "", err
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), err
}

func (r *inventoryRepository) GetStorageByID(ctx context.Context, id *primitive.ObjectID) (*Storage, error) {

	var storage Storage

	filter := bson.M{
		"_id":       id,
		"is_actice": true,
	}

	err := r.storageCollection.FindOne(ctx, filter).Decode(&storage)
	if err != nil {
		return nil, err
	}

	return &storage, nil

}

func (r *inventoryRepository) GetStoragies(ctx context.Context, typeString string) ([]*Storage, error) {

	var storagies []*Storage
	filter := bson.M{
		"type":      typeString,
		"is_actice": true,
	}

	cursor, err := r.storageCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var storage Storage
		if err := cursor.Decode(&storage); err != nil {
			return nil, err
		}
		storagies = append(storagies, &storage)
	}

	return storagies, nil

}

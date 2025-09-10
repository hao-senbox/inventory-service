package storage

import (
	"context"
	"inventory-service/internal/shared/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageRepository interface {
	AddStorage(ctx context.Context, storage *Storage) (string, error)
	GetStoragies(ctx context.Context, typeString string) (map[string][]*Storage, error)
	GetAllStoragies(ctx context.Context) ([]*Storage, error)
	GetStorageByID(ctx context.Context, id *primitive.ObjectID) (*Storage, error)
	GetStorageByShelfID(ctx context.Context, id primitive.ObjectID) ([]*model.Storage, error)
	UpdateStorage(ctx context.Context, id primitive.ObjectID, storage *Storage) error
	DeleteStorage(ctx context.Context, id primitive.ObjectID) error

	// Update total stock
	UpdateTotalStock(ctx context.Context, id primitive.ObjectID, totalStock int) error
}

type storageRepository struct {
	storageCollection *mongo.Collection
}

func NewStorageRepository(storageCollection *mongo.Collection) StorageRepository {
	return &storageRepository{
		storageCollection: storageCollection,
	}
}

func (r *storageRepository) AddStorage(ctx context.Context, storage *Storage) (string, error) {

	result, err := r.storageCollection.InsertOne(ctx, storage)
	if err != nil {
		return "", err
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), err
}

func (r *storageRepository) GetStorageByID(ctx context.Context, id *primitive.ObjectID) (*Storage, error) {

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

func (r *storageRepository) GetStoragies(ctx context.Context, typeString string) (map[string][]*Storage, error) {

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_actice": true}}},
	}

	if typeString != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"type": typeString}}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$type",
			"items": bson.M{"$push": "$$ROOT"},
		}}},
	)

	cursor, err := r.storageCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)


	typeMap := map[string]string{
		"warehouse": "warehouses",
		"building":  "buildings",
		"floor":     "floors",
		"room":      "rooms",
		"shelf":     "shelves",
	}

	results := make(map[string][]*Storage)
	for cursor.Next(ctx) {
		var row struct {
			ID    string     `bson:"_id"`
			Items []*Storage `bson:"items"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}

		if plural, ok := typeMap[row.ID]; ok {
			results[plural] = row.Items
		} else {
			results[row.ID] = row.Items
		}
	}

	allTypes := []string{"warehouses", "buildings", "floors", "rooms", "shelves"}
	for _, t := range allTypes {
		if _, ok := results[t]; !ok {
			results[t] = []*Storage{}
		}
	}

	return results, nil
	
}

func (r *storageRepository) GetAllStoragies(ctx context.Context) ([]*Storage, error) {

	cursor, err := r.storageCollection.Find(ctx, bson.M{"is_actice": true})
	if err != nil {
		return nil, err
	}
	
	defer cursor.Close(ctx)

	var storagies []*Storage
	for cursor.Next(ctx) {
		var storage Storage
		err := cursor.Decode(&storage)
		if err != nil {
			return nil, err
		}
		storagies = append(storagies, &storage)
	}

	return storagies, nil

}

func (r *storageRepository) UpdateStorage(ctx context.Context, id primitive.ObjectID, storage *Storage) error {

	_, err := r.storageCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": storage})
	if err != nil {
		return err
	}

	return nil

}

func (r *storageRepository) DeleteStorage(ctx context.Context, id primitive.ObjectID) error {

	_, err := r.storageCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
	
}

func (r *storageRepository) GetStorageByShelfID(ctx context.Context, id primitive.ObjectID) ([]*model.Storage, error) {

	var storagies []*model.Storage

	filter := bson.M{"shelf_type_id": id}

	cursor, err := r.storageCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var storage *model.Storage
		err := cursor.Decode(&storage)
		if err != nil {
			return nil, err
		}
		storagies = append(storagies, storage)
	}

	return storagies, nil

}

func (r *storageRepository) UpdateTotalStock(ctx context.Context, id primitive.ObjectID, totalStock int) error {

	_, err := r.storageCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"total_stock": totalStock}})
	if err != nil {
		return err
	}

	return nil
}
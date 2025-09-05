package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageRepository interface {
	AddStorage(ctx context.Context, storage *Storage) (string, error)
	GetStoragies(ctx context.Context, typeString string) (map[string][]*Storage, error)
	GetStorageByID(ctx context.Context, id *primitive.ObjectID) (*Storage, error)
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



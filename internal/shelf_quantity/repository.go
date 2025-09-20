package shelfquantity

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShelfQuantityRepository interface {
	CreateShelfQuantity(ctx context.Context, item *ShelfQuantity, userID string) error
	GetShelfQuantitiesByShelfID(ctx context.Context, shelfID primitive.ObjectID) ([]*ShelfQuantity, error)

	DeleteQuantity(ctx context.Context, id primitive.ObjectID) error
}

type shelfQuantityRepository struct {
	ShelfQuantityCollection *mongo.Collection
}

func NewShelfQuantityRepository(collection *mongo.Collection) ShelfQuantityRepository {
	return &shelfQuantityRepository{
		ShelfQuantityCollection: collection,
	}
}

func (r *shelfQuantityRepository) CreateShelfQuantity(ctx context.Context, item *ShelfQuantity, userID string) error {
	_, err := r.ShelfQuantityCollection.InsertOne(ctx, item)
	return err
}

func (r *shelfQuantityRepository) GetShelfQuantitiesByShelfID(ctx context.Context, shelfID primitive.ObjectID) ([]*ShelfQuantity, error) {

	var shelfQuantities []*ShelfQuantity

	cursor, err := r.ShelfQuantityCollection.Find(ctx, bson.M{"shelf_id": shelfID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var shelfQuantity ShelfQuantity
		if err := cursor.Decode(&shelfQuantity); err != nil {
			return nil, err
		}
		shelfQuantities = append(shelfQuantities, &shelfQuantity)
	}

	return shelfQuantities, nil
	
}

func (r *shelfQuantityRepository) DeleteQuantity(ctx context.Context, id primitive.ObjectID) error {

	filter := bson.M{"shelf_id": id}
	
	_, err := r.ShelfQuantityCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
	
}

package inventory

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryRepository interface {
	AddLocation(ctx context.Context, location *Locations) (string, error)
	GetLocations(ctx context.Context, typeString string) ([]*Locations, error)
	GetLocationByID(ctx context.Context, id *primitive.ObjectID) (*Locations, error)
}

type inventoryRepository struct {
	locationCollection *mongo.Collection
}

func NewInventoryRepository(locationCollection *mongo.Collection) InventoryRepository {
	return &inventoryRepository{
		locationCollection: locationCollection,
	}
}

func (r *inventoryRepository) AddLocation(ctx context.Context, location *Locations) (string, error) {
	
	result, err := r.locationCollection.InsertOne(ctx, location)
	if err != nil {
		return "", err
	}

	oid := result.InsertedID.(primitive.ObjectID)
    
	return oid.Hex(), err
}

func (r *inventoryRepository) GetLocationByID(ctx context.Context, id *primitive.ObjectID) (*Locations, error) {

	var location Locations
	fmt.Printf("ID: %v\n", id)
	filter := bson.M{
		"_id": id,
		"is_actice": true,
	}

	err := r.locationCollection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		return nil, err
	}
	
	return &location, nil	

}

func (r *inventoryRepository) GetLocations(ctx context.Context, typeString string) ([]*Locations, error) {

	var locations []*Locations
	filter := bson.M{
		"type": typeString,
		"is_actice": true,
	}

	cursor, err := r.locationCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var location Locations
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}

	return locations, nil
	
}
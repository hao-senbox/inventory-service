package location

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationRepository interface {
	Create(ctx context.Context, location *Location) (string, error)
}

type locationRepository struct {
	locationCollection *mongo.Collection
}

func NewLocationRepository(locationCollection *mongo.Collection) LocationRepository {
	return &locationRepository{
		locationCollection: locationCollection,
	}
}

func (r *locationRepository) Create(ctx context.Context, location *Location) (string, error) {
	
	result, err := r.locationCollection.InsertOne(ctx, location)
	if err != nil {
		return "", err
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), err
	
}		
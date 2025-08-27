package location

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationRepository interface {
	Create(ctx context.Context, location *Location) (string, error)
	GetLocations(ctx context.Context) ([]*Location, error)
	GetLocationByID(ctx context.Context, id primitive.ObjectID) (*Location, error)
	UpdateLocation(ctx context.Context, location *Location) error
	DeleteLocation(ctx context.Context, id primitive.ObjectID) error
	AddProductToLocation(ctx context.Context, locationID primitive.ObjectID, productID string, quantity int) error
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

func (r *locationRepository) GetLocations(ctx context.Context) ([]*Location, error) {

	var locations []*Location

	cursor, err := r.locationCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var location Location
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}

	return locations, nil
	
}

func (r *locationRepository) GetLocationByID(ctx context.Context, id primitive.ObjectID) (*Location, error) {

	var location Location

	if err := r.locationCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (r *locationRepository) UpdateLocation(ctx context.Context, location *Location) error {

	_, err := r.locationCollection.UpdateOne(ctx, bson.M{"_id": location.ID}, bson.M{"$set": location})

	return err

}

func (r *locationRepository) DeleteLocation(ctx context.Context, id primitive.ObjectID) error {

	_, err := r.locationCollection.DeleteOne(ctx, bson.M{"_id": id})
	
	return err
}

func (r *locationRepository) AddProductToLocation(ctx context.Context, locationID primitive.ObjectID, productID string, quantity int) error {
	filter := bson.M{
		"_id":                locationID,
		"products.product_id": productID, 
	}

	update := bson.M{
		"$inc": bson.M{
			"products.$.quantity": quantity,
		},
	}

	result, err := r.locationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		filter = bson.M{
			"_id": locationID,
		}

		update = bson.M{
			"$push": bson.M{
				"products": bson.M{
					"product_id": productID,
					"quantity":   quantity,
				},
			},
		}

		_, err = r.locationCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}

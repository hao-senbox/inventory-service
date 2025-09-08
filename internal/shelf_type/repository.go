package shelftype

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShelfTypeRepository interface {
	CreateShelfType(ctx context.Context, shelfType *ShelfType) (string, error)
	GetShelfTypes(ctx context.Context) ([]*ShelfType, error)
	GetShelfTypeByID(ctx context.Context, id primitive.ObjectID) (*ShelfType, error)
}

type shelfTypeRepository struct {
	collection *mongo.Collection
}

func NewShelfTypeRepository(collection *mongo.Collection) ShelfTypeRepository {
	return &shelfTypeRepository{
		collection: collection,
	}
}

func (s *shelfTypeRepository) CreateShelfType(ctx context.Context, shelfType *ShelfType) (string, error) {

	result, err := s.collection.InsertOne(ctx, shelfType)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("InsertedID is not an ObjectID")
	}

	return oid.Hex(), nil

}

func (s *shelfTypeRepository) GetShelfTypes(ctx context.Context) ([]*ShelfType, error) {

	var shelfTypes []*ShelfType

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var shelfType ShelfType
		err := cursor.Decode(&shelfType)
		if err != nil {
			return nil, err
		}
		shelfTypes = append(shelfTypes, &shelfType)
	}

	return shelfTypes, nil

}

func (s *shelfTypeRepository) GetShelfTypeByID(ctx context.Context, id primitive.ObjectID) (*ShelfType, error) {

	filter := bson.M{"_id": id}

	var shelfType ShelfType

	err := s.collection.FindOne(ctx, filter).Decode(&shelfType)
	if err != nil {
		return nil, err
	}

	return &shelfType, nil

}

package productplacement

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductPlacementRepository interface {
	CreateProductPlacement(ctx context.Context, productPlacement *ProductPlacement) error
	GetByProductAndShelf(ctx context.Context, productID, shelfID primitive.ObjectID) (*ProductPlacement, error)
	ExistsProductPlacement(ctx context.Context, productID, shelfID primitive.ObjectID) (bool, error)
	UpdateProductPlacement(ctx context.Context, productID, shelfID primitive.ObjectID, currentQty int) error
	GetProductPlacementsByProductID(ctx context.Context, productID primitive.ObjectID) ([]*ProductPlacement, error)
	GetProductPlacementsByShelfID(ctx context.Context, shelfID primitive.ObjectID) ([]*ProductPlacement, error)
}

type productPlacementRepository struct {
	collection *mongo.Collection
}

func NewProductPlacementRepository(collection *mongo.Collection) ProductPlacementRepository {
	return &productPlacementRepository{
		collection: collection,
	}
}

func (p *productPlacementRepository) CreateProductPlacement(ctx context.Context, productPlacement *ProductPlacement) error {
	_, err := p.collection.InsertOne(ctx, productPlacement)
	return err
}

func (p *productPlacementRepository) ExistsProductPlacement(ctx context.Context, productID, shelfID primitive.ObjectID) (bool, error) {
	count, err := p.collection.CountDocuments(ctx, bson.M{"product_id": productID, "shelf_id": shelfID})
	return count > 0, err
}

func (p *productPlacementRepository) UpdateProductPlacement(ctx context.Context, productID, shelfID primitive.ObjectID, currentQty int) error {
	_, err := p.collection.UpdateOne(ctx, bson.M{"product_id": productID, "shelf_id": shelfID}, bson.M{"$set": bson.M{"current_qty": currentQty}})
	return err
}

func (p *productPlacementRepository) GetByProductAndShelf(ctx context.Context, productID, shelfID primitive.ObjectID) (*ProductPlacement, error) {

	var placement ProductPlacement

	err := p.collection.FindOne(ctx, bson.M{"product_id": productID, "shelf_id": shelfID}).Decode(&placement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &placement, err
}

func (p *productPlacementRepository) GetProductPlacementsByProductID(ctx context.Context, productID primitive.ObjectID) ([]*ProductPlacement, error) {

	var placements []*ProductPlacement

	cursor, err := p.collection.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var placement ProductPlacement
		if err := cursor.Decode(&placement); err != nil {
			return nil, err
		}
		placements = append(placements, &placement)
	}

	return placements, nil
}

func (p *productPlacementRepository) GetProductPlacementsByShelfID(ctx context.Context, shelfID primitive.ObjectID) ([]*ProductPlacement, error) {

	var placements []*ProductPlacement

	cursor, err := p.collection.Find(ctx, bson.M{"shelf_id": shelfID})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var placement ProductPlacement
		if err := cursor.Decode(&placement); err != nil {
			return nil, err
		}
		placements = append(placements, &placement)
	}

	return placements, nil
}
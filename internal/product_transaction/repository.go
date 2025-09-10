package producttransaction

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProductTransactionRepository interface {
	CreateProductTransaction(ctx context.Context, data *ProductTransaction) error
}

type productTransactionRepository struct {
	collection *mongo.Collection
}

func NewProductTransactionRepository(collection *mongo.Collection) ProductTransactionRepository {
	return &productTransactionRepository{
		collection: collection,
	}
}

func (p *productTransactionRepository) CreateProductTransaction(ctx context.Context, data *ProductTransaction) error {

	_, err := p.collection.InsertOne(ctx, data)

	if err != nil {
		return err
	}

	return nil
	
}
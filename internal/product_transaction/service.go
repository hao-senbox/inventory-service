package producttransaction

import (
	"context"
	"fmt"
	productplacement "inventory-service/internal/product_placement"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductTransactionService interface {
	CreateProductTransaction(ctx context.Context, req *CreateProductTransactionRequest, userID string) (string, error)
}

type productTransactionService struct {
	ProductTransactionRepository ProductTransactionRepository
	ProductPlacementService      productplacement.ProductPlacementService
	mongoClient                  *mongo.Client
}

func NewProductTransactionService(
	productTransactionRepository ProductTransactionRepository,
	productPlacementService productplacement.ProductPlacementService,
	mongoClient *mongo.Client,
) ProductTransactionService {
	return &productTransactionService{
		ProductTransactionRepository: productTransactionRepository,
		ProductPlacementService:      productPlacementService,
		mongoClient:                  mongoClient,
	}
}

func (s *productTransactionService) CreateProductTransaction(ctx context.Context, req *CreateProductTransactionRequest, userID string) (string, error) {

	if userID == "" {
		return "", fmt.Errorf("user_id is required")
	}

	if req.ProductID == "" {
		return "", fmt.Errorf("product_id is required")
	}

	if req.ShelfID == "" {
		return "", fmt.Errorf("shelf_id is required")
	}

	if req.Quantity <= 0 {
		return "", fmt.Errorf("quantity must be greater than 0")
	}

	if req.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	objProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return "", fmt.Errorf("invalid product id: %v", err)
	}

	objShelfID, err := primitive.ObjectIDFromHex(req.ShelfID)
	if err != nil {
		return "", fmt.Errorf("invalid shelf id: %v", err)
	}

	ID := primitive.NewObjectID()
	productTransaction := &ProductTransaction{
		ID:        ID,
		ProductID: objProductID,
		ShelfID:   objShelfID,
		Quantity:  req.Quantity,
		Action:    req.Action,
		ActionBy:  userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(ctx)

	callback := func(sc mongo.SessionContext) (interface{}, error) {

		if err := s.ProductTransactionRepository.CreateProductTransaction(sc, productTransaction); err != nil {
			return nil, err
		}

		if req.Action == "IN" {
			placementReq := &productplacement.CreateProductPlacementRequest{
				ProductID:  req.ProductID,
				ShelfID:    req.ShelfID,
				CurrentQty: req.Quantity,
			}
			if err := s.ProductPlacementService.CreateProductPlacement(sc, placementReq); err != nil {
				return nil, err
			}
		} else if req.Action == "OUT" {
			placementReq := &productplacement.UpdateProductPlacementRequest{
				ProductID:  req.ProductID,
				ShelfID:    req.ShelfID,
				CurrentQty: req.Quantity,
			}
			if err := s.ProductPlacementService.UpdateProductPlacement(sc, placementReq); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return "", err
	}

	return ID.Hex(), nil

}

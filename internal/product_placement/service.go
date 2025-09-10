package productplacement

import (
	"context"
	"fmt"
	"inventory-service/internal/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductPlacementService interface {
	GetProductPlacementsByShelfID(ctx context.Context, shelfId string) ([]*ProductPlacement, error)
	GetProductPlacementsByProductID(ctx context.Context, productId string) ([]*ProductPlacement, error)
	CreateProductPlacement(ctx context.Context, req *CreateProductPlacementRequest) error
	UpdateProductPlacement(ctx context.Context, req *UpdateProductPlacementRequest) error
}

type productPlacementService struct {
	repository        ProductPlacementRepository
	storageRepository storage.StorageRepository
}

func NewProductPlacementService(repository ProductPlacementRepository, storageRepository storage.StorageRepository) ProductPlacementService {
	return &productPlacementService{
		repository:        repository,
		storageRepository: storageRepository,
	}
}

func (p *productPlacementService) CreateProductPlacement(ctx context.Context, req *CreateProductPlacementRequest) error {

	sc, ok := ctx.(mongo.SessionContext)
	if !ok {
		return fmt.Errorf("context is not a session context")
	}

	if req.ProductID == "" {
		return fmt.Errorf("product_id is required")
	}
	if req.ShelfID == "" {
		return fmt.Errorf("shelf_id is required")
	}
	if req.CurrentQty <= 0 {
		return fmt.Errorf("current_qty must be greater than 0")
	}

	objProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	objShelfID, err := primitive.ObjectIDFromHex(req.ShelfID)
	if err != nil {
		return fmt.Errorf("invalid shelf id: %v", err)
	}

	placement, err := p.repository.GetByProductAndShelf(sc, objProductID, objShelfID)
	if err != nil {
		return err
	}

	storage, err := p.storageRepository.GetStorageByID(sc, &objShelfID)
	if err != nil {
		return err
	}
	if storage.TotalStock == nil || *storage.TotalStock < req.CurrentQty {
		return fmt.Errorf("not enough stock capacity")
	}
	
	if placement != nil {
		err = p.repository.UpdateProductPlacement(sc, objProductID, objShelfID, placement.CurrentQty+req.CurrentQty)
		if err != nil {
			return err
		}
	} else {
		data := &ProductPlacement{
			ID:          primitive.NewObjectID(),
			ProductID:   objProductID,
			ShelfID:     objShelfID,
			CurrentQty:  req.CurrentQty,
			Path:        storage.Path,
			AncestorIDs: storage.AncestorIDs,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := p.repository.CreateProductPlacement(sc, data); err != nil {
			return err
		}
	}

	return p.storageRepository.UpdateTotalStock(sc, objShelfID, *storage.TotalStock-req.CurrentQty)
}

func (p *productPlacementService) UpdateProductPlacement(ctx context.Context, req *UpdateProductPlacementRequest) error {

	sc, ok := ctx.(mongo.SessionContext)
	if !ok {
		return fmt.Errorf("context is not a session context")
	}

	if req.ProductID == "" {
		return fmt.Errorf("product_id is required")
	}

	if req.ShelfID == "" {
		return fmt.Errorf("shelf_id is required")
	}

	if req.CurrentQty <= 0 {
		return fmt.Errorf("current_qty must be greater than 0")
	}

	objProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	objShelfID, err := primitive.ObjectIDFromHex(req.ShelfID)
	if err != nil {
		return fmt.Errorf("invalid shelf id: %v", err)
	}

	placement, err := p.repository.GetByProductAndShelf(sc, objProductID, objShelfID)
	if err != nil {
		return err
	}

	if placement == nil {
		return fmt.Errorf("product placement not found")
	}

	newQty := placement.CurrentQty - req.CurrentQty
	if newQty < 0 {
		return fmt.Errorf("not enough stock to OUT")
	}

	err = p.repository.UpdateProductPlacement(sc, objProductID, objShelfID, newQty)
	if err != nil {
		return err
	}

	storage, err := p.storageRepository.GetStorageByID(sc, &objShelfID)
	if err != nil {
		return err
	}

	err = p.storageRepository.UpdateTotalStock(sc, objShelfID, *storage.TotalStock+req.CurrentQty)
	if err != nil {
		return err
	}

	return nil
}

func (p *productPlacementService) GetProductPlacementsByProductID(ctx context.Context, productID string) ([]*ProductPlacement, error) {
	
	if productID == "" {
		return nil, fmt.Errorf("product_id is required")
	}

	objProductID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %v", err)
	}

	return p.repository.GetProductPlacementsByProductID(ctx, objProductID)
}

func (p *productPlacementService) GetProductPlacementsByShelfID(ctx context.Context, shelfId string) ([]*ProductPlacement, error) {
	
	if shelfId == "" {
		return nil, fmt.Errorf("shelf_id is required")
	}

	objShelfID, err := primitive.ObjectIDFromHex(shelfId)
	if err != nil {
		return nil, fmt.Errorf("invalid shelf id: %v", err)
	}

	return p.repository.GetProductPlacementsByShelfID(ctx, objShelfID)

}

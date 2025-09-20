package shelfquantity

import (
	"context"
	"fmt"
	"inventory-service/internal/shared/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShelfQuantityService interface {
	CreateShelfQuantity(ctx context.Context, req []CreateShelfQuantityRequest, userID string) error
	GetShelfQuantitiesByShelfID(ctx context.Context, shelfID string) ([]*ShelfQuantity, error)
}

type shelfQuantityService struct {
	ShelfQuantityRepository ShelfQuantityRepository
	StorageRepository       ports.Storage
}

func NewShelfQuantityService(shelfQuantityRepository ShelfQuantityRepository,
	storageRepository ports.Storage) ShelfQuantityService {
	return &shelfQuantityService{
		ShelfQuantityRepository: shelfQuantityRepository,
		StorageRepository:       storageRepository,
	}
}

func (s *shelfQuantityService) CreateShelfQuantity(ctx context.Context, req []CreateShelfQuantityRequest, userID string) error {

	shelfID := req[0].ShelfID

	shelfObjectID, err := primitive.ObjectIDFromHex(shelfID)
	if err != nil {
		return fmt.Errorf("invalid shelf id: %v", err)
	}

	result, err := s.StorageRepository.GetStorageByID(ctx, &shelfObjectID)
	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("storage not found")
	}

	for _, item := range req {

		if item.Code == "" {
			return fmt.Errorf("code is required")
		}

		if item.Note == "" {
			return fmt.Errorf("note is required")
		}

		objectID, err := primitive.ObjectIDFromHex(item.ShelfID)
		if err != nil {
			return fmt.Errorf("invalid shelf id: %v", err)
		}

		qrCode := fmt.Sprintf("SENBOX.ORG[INVENTORY]:%s", item.Code)

		data := &ShelfQuantity{
			ID:        primitive.NewObjectID(),
			ShelfID:   objectID,
			Code:      item.Code,
			QRCode:    qrCode,
			Note:      item.Note,
			CreatedBy: userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = s.ShelfQuantityRepository.CreateShelfQuantity(ctx, data, userID)
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *shelfQuantityService) GetShelfQuantitiesByShelfID(ctx context.Context, shelfID string) ([]*ShelfQuantity, error) {

	objectID, err := primitive.ObjectIDFromHex(shelfID)
	if err != nil {
		return nil, fmt.Errorf("invalid shelf id: %v", err)
	}

	return s.ShelfQuantityRepository.GetShelfQuantitiesByShelfID(ctx, objectID)
}

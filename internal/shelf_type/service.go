package shelftype

import (
	"context"
	"fmt"
	"inventory-service/internal/shared/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShelfTypeService interface {
	CreateShelfType(ctx context.Context, req *CreateShelfTypeRequest) (string, error)
	GetShelfTypes(ctx context.Context) ([]*ShelfType, error)
	GetShelfTypeByID(ctx context.Context, id string) (*ShelfType, error)
	UpdateShelfType(ctx context.Context, id string, req *UpdateShelfTypeRequest) error
	DeleteShelfType(ctx context.Context, id string) error
}

type shelfTypeService struct {
	ShelfTypeRepo ShelfTypeRepository
	StorageRepo   ports.Storage
}

func NewShelfTypeService(shelfTypeRepo ShelfTypeRepository, storageRepo ports.Storage) ShelfTypeService {
	return &shelfTypeService{
		ShelfTypeRepo: shelfTypeRepo,
		StorageRepo:   storageRepo,
	}
}

func (s *shelfTypeService) CreateShelfType(ctx context.Context, req *CreateShelfTypeRequest) (string, error) {

	var stock *int

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	
	if req.Slot != nil && req.Level != nil {
		val := (*req.Level) * (*req.Slot)
		stock = &val
	} else {
		stock = nil
	}
	
	shelfType := &ShelfType{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Note:      req.Note,
		Slot:      req.Slot,
		Level:     req.Level,
		Stock:     stock,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.ShelfTypeRepo.CreateShelfType(ctx, shelfType)
	if err != nil {
		return "", err
	}

	return id, nil

}

func (s *shelfTypeService) GetShelfTypes(ctx context.Context) ([]*ShelfType, error) {
	return s.ShelfTypeRepo.GetShelfTypes(ctx)
}

func (s *shelfTypeService) GetShelfTypeByID(ctx context.Context, id string) (*ShelfType, error) {

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.ShelfTypeRepo.GetShelfTypeByID(ctx, objectID)

}

func (s *shelfTypeService) UpdateShelfType(ctx context.Context, id string, req *UpdateShelfTypeRequest) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	shelfType, err := s.ShelfTypeRepo.GetShelfTypeByID(ctx, objectID)
	if err != nil {
		return err
	}

	if shelfType == nil {
		return fmt.Errorf("shelf type not found")
	}

	if req.Name != "" {
		shelfType.Name = req.Name
	}

	if req.Note != nil {
		shelfType.Note = req.Note
	}

	if req.Slot != nil {
		shelfType.Slot = req.Slot
	}

	if req.Level != nil {
		shelfType.Level = req.Level
	}

	shelfType.UpdatedAt = time.Now()

	return s.ShelfTypeRepo.UpdateShelfType(ctx, objectID, shelfType)

}

func (s *shelfTypeService) DeleteShelfType(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	storagies, err := s.StorageRepo.GetStorageByShelfID(ctx, objectID)
	if err != nil {
		return err
	}

	if storagies == nil {
		return fmt.Errorf("storage type not found")
	}

	for _, storage := range storagies {
		err := s.StorageRepo.DeleteStorage(ctx, storage.ID)
		if err != nil {
			return err
		}
	}

	err = s.ShelfTypeRepo.DeleteShelfType(ctx, objectID)
	if err != nil {
		return err
	}

	return nil
}

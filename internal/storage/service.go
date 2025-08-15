package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StorageService interface {
	CreateStorage(ctx context.Context, req *CreateStorageRequest, userID string) (string, error)
	GetStoragies(ctx context.Context, typeString string) ([]*Storage, error)
	GetStorageByID(ctx context.Context, id string) (*Storage, error)
}

type storageService struct {
	repository StorageRepository
}

func NewStorageService(repository StorageRepository) StorageService {
	return &storageService{
		repository: repository,
	}
}

func (s *storageService) CreateStorage(ctx context.Context, req *CreateStorageRequest, userID string) (string, error) {

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if req.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	if req.ImageMain == "" {
		return "", fmt.Errorf("image_main is required")
	}

	if req.ImageMap == "" {
		return "", fmt.Errorf("image_map is required")
	}

	storage := &Storage{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Type:        req.Type,
		Description: &req.Description,
		ImageMain:   req.ImageMain,
		ImageMap:    req.ImageMap,
		IsActice:    true,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.buildLocationHierarchy(ctx, storage); err != nil {
		return "", err
	}

	storageID, err := s.repository.AddStorage(ctx, storage)
	if err != nil {
		return "", err
	}

	return storageID, nil
}

func (s *storageService) buildLocationHierarchy(ctx context.Context, storage *Storage) error {

	if storage.ParentID == nil {
		storage.Level = 0
		storage.Path = "/" + storage.Name
		storage.AncestorIDs = []primitive.ObjectID{}
		return nil
	}

	parent, err := s.repository.GetStorageByID(ctx, storage.ParentID)
	if err != nil {
		return err
	}

	storage.Level = parent.Level + 1
	storage.Path = parent.Path + "/" + storage.Name
	storage.AncestorIDs = append(parent.AncestorIDs, parent.ID)

	return nil
}

func (s *storageService) GetStoragies(ctx context.Context, typeString string) ([]*Storage, error) {
	return s.repository.GetStoragies(ctx, typeString)
}

func (s *storageService) GetStorageByID(ctx context.Context, id string) (*Storage, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repository.GetStorageByID(ctx, &objectID)
	
}
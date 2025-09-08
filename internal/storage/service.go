package storage

import (
	"context"
	"fmt"
	shelftype "inventory-service/internal/shelf_type"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StorageService interface {
	CreateStorage(ctx context.Context, req *CreateStorageRequest, userID string) (string, error)
	GetStoragies(ctx context.Context, typeString string) (map[string][]*Storage, error)
	GetStorageByID(ctx context.Context, id string) (*Storage, error)
	GetStorageTree(ctx context.Context) ([]*StorageNodeResponse, error)
}

type storageService struct {
	repository          StorageRepository
	shelfTypeRepository shelftype.ShelfTypeRepository
}

func NewStorageService(repository StorageRepository, shelfTypeRepository shelftype.ShelfTypeRepository) StorageService {
	return &storageService{
		repository:          repository,
		shelfTypeRepository: shelfTypeRepository,
	}
}

func (s *storageService) CreateStorage(ctx context.Context, req *CreateStorageRequest, userID string) (string, error) {

	var parentID *primitive.ObjectID
	var shelfTypeID *primitive.ObjectID

	if req.ParentID != nil {
		parentIDConvert, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return "", err
		}
		parentID = &parentIDConvert
	} else {
		parentID = nil
	}

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if req.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	var storage *Storage
	ID := primitive.NewObjectID()
	QRCocde := fmt.Sprintf("SENBOX.ORG[STORAGE]:%s", ID.Hex())

	if req.ShelfTypeID != nil {

		ShelfTypeIDConvert, err := primitive.ObjectIDFromHex(*req.ShelfTypeID)
		if err != nil {
			return "", err
		}

		shelfTypeID = &ShelfTypeIDConvert

		shelfType, err := s.shelfTypeRepository.GetShelfTypeByID(ctx, *shelfTypeID)
		if err != nil {
			return "", err
		}

		if shelfType == nil {
			return "", fmt.Errorf("shelf type not found")
		}

		storage = &Storage{
			ID:          ID,
			Name:        req.Name,
			Type:        req.Type,
			QRCode:      QRCocde,
			Description: &req.Description,
			ImageMain:   req.ImageMain,
			ImageMap:    req.ImageMap,
			ParentID:    parentID,
			ShelfTypeID: shelfTypeID,
			Slots:       shelfType.Slot,
			Levels:      shelfType.Level,
			IsActice:    true,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	} else {
		storage = &Storage{
			ID:          ID,
			Name:        req.Name,
			Type:        req.Type,
			QRCode:      QRCocde,
			Description: &req.Description,
			ImageMain:   req.ImageMain,
			ImageMap:    req.ImageMap,
			ParentID:    parentID,
			IsActice:    true,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
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

func (s *storageService) GetStoragies(ctx context.Context, typeString string) (map[string][]*Storage, error) {
	return s.repository.GetStoragies(ctx, typeString)
}

func (s *storageService) GetStorageByID(ctx context.Context, id string) (*Storage, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repository.GetStorageByID(ctx, &objectID)

}

func (s *storageService) GetStorageTree(ctx context.Context) ([]*StorageNodeResponse, error) {
	
	storagies, err := s.repository.GetAllStoragies(ctx)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[string]*StorageNodeResponse)
	for _, storage := range storagies {
		nodeMap[storage.ID.Hex()] = &StorageNodeResponse{
			Storage:  *storage,
			Children: []*StorageNodeResponse{},
		}
	}

	var roots []*StorageNodeResponse
	for _, storage := range storagies {
		if storage.ParentID == nil {
			roots = append(roots, nodeMap[storage.ID.Hex()])
		} else {
			nodeMap[storage.ParentID.Hex()].Children = append(nodeMap[storage.ParentID.Hex()].Children, nodeMap[storage.ID.Hex()])
		}
	}

	return roots, nil

}
package storage

import (
	"context"
	"fmt"
	shelftype "inventory-service/internal/shelf_type"
	"inventory-service/pkg/uploader"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StorageService interface {
	CreateStorage(ctx context.Context, req *CreateStorageRequest, userID string) (string, error)
	GetStoragies(ctx context.Context, typeString string) (map[string][]*Storage, error)
	GetStorageByID(ctx context.Context, id string) (*Storage, error)
	GetStorageTree(ctx context.Context) ([]*StorageNodeResponse, error)
	UpdateStorage(ctx context.Context, id string, req *UpdateStorageRequest) error
	DeleteStorage(ctx context.Context, id string) error
}

type storageService struct {
	repository          StorageRepository
	shelfTypeRepository shelftype.ShelfTypeRepository
	ImageService        uploader.ImageService
}

func NewStorageService(repository StorageRepository, shelfTypeRepository shelftype.ShelfTypeRepository, imageService uploader.ImageService) StorageService {
	return &storageService{
		repository:          repository,
		shelfTypeRepository: shelfTypeRepository,
		ImageService:        imageService,
	}
}

func (s *storageService) CreateStorage(ctx context.Context, req *CreateStorageRequest, userID string) (string, error) {

	var parentID *primitive.ObjectID
	var shelfTypeID *primitive.ObjectID

	if req.ParentID != nil {
		parentIDConvert, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return "", fmt.Errorf("invalid parent id: %w", err)
		}
		parentID = &parentIDConvert
	}

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	if req.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	ID := primitive.NewObjectID()
	qrCode := fmt.Sprintf("SENBOX.ORG[STORAGE]:%s", ID.Hex())

	var storage *Storage

	if req.ShelfTypeID != nil {
		ShelfTypeIDConvert, err := primitive.ObjectIDFromHex(*req.ShelfTypeID)
		if err != nil {
			return "", fmt.Errorf("invalid shelf type id: %w", err)
		}
		shelfTypeID = &ShelfTypeIDConvert

		shelfType, err := s.shelfTypeRepository.GetShelfTypeByID(ctx, *shelfTypeID)
		if err != nil {
			return "", err
		}
		if shelfType == nil {
			return "", fmt.Errorf("shelf type not found")
		}

		var totalStock *int
		if shelfType.Slot != nil && shelfType.Level != nil {
			val := (*shelfType.Slot) * (*shelfType.Level)
			totalStock = &val
		}

		storage = &Storage{
			ID:          ID,
			Name:        req.Name,
			Type:        req.Type,
			QRCode:      qrCode,
			Description: &req.Description,
			ImageMain:   req.ImageMain,
			ImageMap:    req.ImageMap,
			ParentID:    parentID,
			ShelfTypeID: shelfTypeID,
			Slots:       shelfType.Slot,
			Levels:      shelfType.Level,
			TotalStock:  totalStock,
			IsActive:    true,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	} else {
		storage = &Storage{
			ID:          ID,
			Name:        req.Name,
			Type:        req.Type,
			QRCode:      qrCode,
			Description: &req.Description,
			ImageMain:   req.ImageMain,
			ImageMap:    req.ImageMap,
			ParentID:    parentID,
			IsActive:    true,
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
		var imageMainUrl string
		if storage.ImageMain != nil && *storage.ImageMain != "" {
			urlImage, err := s.ImageService.GetImageKey(ctx, *storage.ImageMain)
			if err == nil && urlImage != nil {
				imageMainUrl = urlImage.Url
			} else {
				imageMainUrl = ""
			}
		}

		var imageMapUrl string
		if storage.ImageMap != nil && *storage.ImageMap != "" {
			urlImage, err := s.ImageService.GetImageKey(ctx, *storage.ImageMap)
			if err == nil && urlImage != nil {
				imageMapUrl = urlImage.Url
			} else {
				imageMapUrl = ""
			}
		}

		nodeMap[storage.ID.Hex()] = &StorageNodeResponse{
			Storage:      *storage,
			ImageMainUrl: imageMainUrl,
			ImageMapUrl:  imageMapUrl,
			Children:     []*StorageNodeResponse{},
		}
	}

	var roots []*StorageNodeResponse
	for _, storage := range storagies {
		if storage.ParentID == nil {
			roots = append(roots, nodeMap[storage.ID.Hex()])
		} else {
			parentNode, ok := nodeMap[storage.ParentID.Hex()]
			if ok {
				parentNode.Children = append(parentNode.Children, nodeMap[storage.ID.Hex()])
			}
		}
	}

	return roots, nil
}

func (s *storageService) UpdateStorage(ctx context.Context, id string, req *UpdateStorageRequest) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	storage, err := s.repository.GetStorageByID(ctx, &objectID)
	if err != nil {
		return err
	}

	if storage == nil {
		return fmt.Errorf("storage not found")
	}

	if req.Name != "" {
		storage.Name = req.Name
	}

	if req.Type != "" {
		storage.Type = req.Type
	}

	if req.Description != nil {
		storage.Description = req.Description
	}

	if req.ImageMain != nil {
		if storage.ImageMain != nil {
			if err := s.ImageService.DeleteImageKey(ctx, *storage.ImageMain); err != nil {
				return err
			}
		}
		storage.ImageMain = req.ImageMain
	}

	if req.ImageMap != nil {
		if storage.ImageMap != nil {
			if err := s.ImageService.DeleteImageKey(ctx, *storage.ImageMap); err != nil {
				return err
			}
		}
		storage.ImageMap = req.ImageMap
	}

	if req.ParentID != nil {
		if *req.ParentID == "" {
			storage.ParentID = nil
		} else {
			parentID, err := primitive.ObjectIDFromHex(*req.ParentID)
			if err != nil {
				return fmt.Errorf("invalid parent id: %w", err)
			}
			storage.ParentID = &parentID
		}
		if err := s.buildLocationHierarchy(ctx, storage); err != nil {
			return err
		}
	}

	if req.ShelfTypeID != nil {
		shelfTypeID, err := primitive.ObjectIDFromHex(*req.ShelfTypeID)
		if err != nil {
			return fmt.Errorf("invalid shelf type id: %w", err)
		}
		storage.ShelfTypeID = &shelfTypeID

		shelfType, err := s.shelfTypeRepository.GetShelfTypeByID(ctx, shelfTypeID)
		if err != nil {
			return err
		}
		if shelfType == nil {
			return fmt.Errorf("shelf type not found")
		}

		var totalStock *int
		if shelfType.Slot != nil && shelfType.Level != nil {
			val := (*shelfType.Slot) * (*shelfType.Level)
			totalStock = &val
		}

		storage.Slots = shelfType.Slot
		storage.Levels = shelfType.Level
		storage.TotalStock = totalStock
	}

	storage.UpdatedAt = time.Now()

	if err := s.repository.UpdateStorage(ctx, objectID, storage); err != nil {
		return err
	}

	return nil
}

func (s *storageService) DeleteStorage(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	storage, err := s.repository.GetStorageByID(ctx, &objectID)
	if err != nil {
		return err
	}

	if storage == nil {
		return fmt.Errorf("storage not found")
	}

	if storage.ImageMain != nil {
		if err := s.ImageService.DeleteImageKey(ctx, *storage.ImageMain); err != nil {
			return err
		}
	}
	if storage.ImageMap != nil {
		if err := s.ImageService.DeleteImageKey(ctx, *storage.ImageMap); err != nil {
			return err
		}
	}

	return s.repository.DeleteStorage(ctx, objectID)

}

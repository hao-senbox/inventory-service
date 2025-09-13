package shelftype

import (
	"context"
	"fmt"
	"inventory-service/internal/shared/ports"
	"inventory-service/pkg/uploader"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShelfTypeService interface {
	CreateShelfType(ctx context.Context, req *CreateShelfTypeRequest) (string, error)
	GetShelfTypes(ctx context.Context) ([]*ShelfTypeResponse, error)
	GetShelfTypeByID(ctx context.Context, id string) (*ShelfTypeResponse, error)
	UpdateShelfType(ctx context.Context, id string, req *UpdateShelfTypeRequest) error
	DeleteShelfType(ctx context.Context, id string) error
}

type shelfTypeService struct {
	ShelfTypeRepo ShelfTypeRepository
	StorageRepo   ports.Storage
	ImageService  uploader.ImageService
}

func NewShelfTypeService(shelfTypeRepo ShelfTypeRepository, storageRepo ports.Storage, imageService uploader.ImageService) ShelfTypeService {
	return &shelfTypeService{
		ShelfTypeRepo: shelfTypeRepo,
		StorageRepo:   storageRepo,
		ImageService:  imageService,
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
		ImageKey:  req.ImageKey,
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

func (s *shelfTypeService) GetShelfTypes(ctx context.Context) ([]*ShelfTypeResponse, error) {

	shelves, err := s.ShelfTypeRepo.GetShelfTypes(ctx)
	if err != nil {
		return nil, err
	}

	var shelfTypes []*ShelfTypeResponse

	for _, shelf := range shelves {

		var imageUrl string

		urlImage, err := s.ImageService.GetImageKey(ctx, shelf.ImageKey)
		if err != nil {
			return nil, err
		}
		fmt.Printf("image url: %v\n", urlImage)
		if urlImage != nil {
			imageUrl = urlImage.Url
		} else {
			imageUrl = ""
		}

		shelfTypes = append(shelfTypes, &ShelfTypeResponse{
			ID:        shelf.ID,
			ImageUrl:  imageUrl,
			Name:      shelf.Name,
			Note:      shelf.Note,
			Slot:      shelf.Slot,
			Level:     shelf.Level,
			Stock:     shelf.Stock,
			CreatedAt: shelf.CreatedAt,
			UpdatedAt: shelf.UpdatedAt,
		})
	}

	return shelfTypes, nil
}

func (s *shelfTypeService) GetShelfTypeByID(ctx context.Context, id string) (*ShelfTypeResponse, error) {

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	shelf, err := s.ShelfTypeRepo.GetShelfTypeByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if shelf == nil {
		return nil, fmt.Errorf("shelf type not found")
	}

	var imageUrl string

	urlImage, err := s.ImageService.GetImageKey(ctx, shelf.ImageKey)
	if err != nil {
		return nil, err
	}

	if urlImage != nil {
		imageUrl = urlImage.Url
	} else {
		imageUrl = ""
	}

	shelfType := &ShelfTypeResponse{
		ID:        shelf.ID,
		ImageUrl:  imageUrl,
		Name:      shelf.Name,
		Note:      shelf.Note,
		Slot:      shelf.Slot,
		Level:     shelf.Level,
		Stock:     shelf.Stock,
		CreatedAt: shelf.CreatedAt,
		UpdatedAt: shelf.UpdatedAt,
	}

	return shelfType, nil

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

	if req.ImageKey != "" {
		err := s.ImageService.DeleteImageKey(ctx, shelfType.ImageKey)
		if err != nil {
			return err
		} 

		shelfType.ImageKey = req.ImageKey
	}

	if req.Name != "" {
		shelfType.Name = req.Name
	}

	if req.Note != nil {
		shelfType.Note = req.Note
	}

	needRecalcStock := false

	if req.Slot != nil {
		shelfType.Slot = req.Slot
		needRecalcStock = true
	}

	if req.Level != nil {
		shelfType.Level = req.Level
		needRecalcStock = true
	}

	if needRecalcStock && shelfType.Slot != nil && shelfType.Level != nil {
		stock := (*shelfType.Slot) * (*shelfType.Level)
		shelfType.Stock = &stock
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

	shelf, err := s.ShelfTypeRepo.GetShelfTypeByID(ctx, objectID)
	if err != nil {
		return err
	}

	if shelf == nil {
		return fmt.Errorf("shelf type not found")
	}

	err = s.ImageService.DeleteImageKey(ctx, shelf.ImageKey)
	if err != nil {
		return err
	}

	storagies, err := s.StorageRepo.GetStorageByShelfID(ctx, objectID)
	if err != nil {
		log.Println(err)
	}

	if storagies == nil {
		log.Println("storage not found")
	}

	// for _, storage := range storagies {
	// 	err := s.StorageRepo.DeleteStorage(ctx, storage.ID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	err = s.ShelfTypeRepo.DeleteShelfType(ctx, objectID)
	if err != nil {
		return err
	}

	return nil
}

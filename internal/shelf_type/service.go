package shelftype

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShelfTypeService interface {
	CreateShelfType(ctx context.Context, req *CreateShelfTypeRequest) (string, error)
	GetShelfTypes(ctx context.Context) ([]*ShelfType, error)
	GetShelfTypeByID(ctx context.Context, id string) (*ShelfType, error)
}

type shelfTypeService struct {
	ShelfTypeRepo ShelfTypeRepository
}

func NewShelfTypeService(shelfTypeRepo ShelfTypeRepository) ShelfTypeService {
	return &shelfTypeService{
		ShelfTypeRepo: shelfTypeRepo,
	}
}

func (s *shelfTypeService) CreateShelfType(ctx context.Context, req *CreateShelfTypeRequest) (string, error) {

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	shelfType := &ShelfType{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Note:      req.Note,
		Slot:      req.Slot,
		Level:     req.Level,
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

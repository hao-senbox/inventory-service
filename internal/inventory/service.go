package inventory

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryService interface {
	CreateLocation(ctx context.Context, req *CreateLocationRequest, userID string) (string, error)
	GetLocations(ctx context.Context, typeString string) ([]*Locations, error)
	GetLocationByID(ctx context.Context, id string) (*Locations, error)
}

type inventoryService struct {
	repository InventoryRepository
}

func NewInventoryService(repository InventoryRepository) InventoryService {
	return &inventoryService{
		repository: repository,
	}
}

func (s *inventoryService) CreateLocation(ctx context.Context, req *CreateLocationRequest, userID string) (string, error) {

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

	var parendID *primitive.ObjectID
	if req.ParentID != nil {
		parentID, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return "", err
		}
		parendID = &parentID
	} else {
		parendID = nil
	}

	location := &Locations{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Type:        req.Type,
		Description: &req.Description,
		ImageMain:   req.ImageMain,
		ImageMap:    req.ImageMap,
		ParentID:    parendID,
		IsActice:    true,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.buildLocationHierarchy(ctx, location); err != nil {
		return "", err
	}

	locationID, err := s.repository.AddLocation(ctx, location)
	if err != nil {
		return "", err
	}

	return locationID, nil
}

func (s *inventoryService) buildLocationHierarchy(ctx context.Context, location *Locations) error {

	if location.ParentID == nil {
		location.Level = 0
		location.Path = "/" + location.Name
		location.AncestorIDs = []primitive.ObjectID{}
		return nil
	}

	parent, err := s.repository.GetLocationByID(ctx, location.ParentID)
	if err != nil {
		return err
	}

	location.Level = parent.Level + 1
	location.Path = parent.Path + "/" + location.Name
	location.AncestorIDs = append(parent.AncestorIDs, parent.ID)

	return nil
}

func (s *inventoryService) GetLocations(ctx context.Context, typeString string) ([]*Locations, error) {
	return s.repository.GetLocations(ctx, typeString)
}

func (s *inventoryService) GetLocationByID(ctx context.Context, id string) (*Locations, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repository.GetLocationByID(ctx, &objectID)
	
}
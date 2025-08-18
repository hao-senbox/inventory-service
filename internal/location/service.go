package location

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationService interface {
	CreateLocation(ctx context.Context, req *CreateLocationRequest, userID string) (string, error)
}

type locationService struct {
	repository LocationRepository
}

func NewLocationService(repository LocationRepository) LocationService {
	return &locationService{
		repository: repository,
	}
}

func (s *locationService) CreateLocation(ctx context.Context, req *CreateLocationRequest, userID string) (string, error) {

	var objShelfID []primitive.ObjectID

	if userID == "" {
		return "", fmt.Errorf("user_id is required")
	}

	if req.LocationName == "" {
		return "", fmt.Errorf("location_name is required")
	}

	if req.LocationAddress == "" {
		return "", fmt.Errorf("location_address is required")
	}

	if req.LocationImage == "" {
		return "", fmt.Errorf("location_image is required")
	}

	if req.LocationMap == "" {
		return "", fmt.Errorf("location_map is required")
	}

	if req.WareHouseID == "" {
		return "", fmt.Errorf("warehouse_id is required")
	}

	wareHouseObjID, err := primitive.ObjectIDFromHex(req.WareHouseID)
	if err != nil {
		return "", fmt.Errorf("warehouse_id is invalid")
	}

	if req.BuildingID == "" {
		return "", fmt.Errorf("building_id is required")
	}

	buildingObjID, err := primitive.ObjectIDFromHex(req.BuildingID)
	if err != nil {
		return "", fmt.Errorf("building_id is invalid")
	}

	if req.FloorID == "" {
		return "", fmt.Errorf("floor_id is required")
	}

	floorObjID, err := primitive.ObjectIDFromHex(req.FloorID)
	if err != nil {
		return "", fmt.Errorf("floor_id is invalid")
	}

	if req.RoomID == "" {
		return "", fmt.Errorf("room_id is required")
	}

	roomObjID, err := primitive.ObjectIDFromHex(req.RoomID)
	if err != nil {
		return "", fmt.Errorf("room_id is invalid")
	}

	if len(req.ShelfID) == 0 {
		return "", fmt.Errorf("shelf_id is required")
	}

	for _, shelfID := range req.ShelfID {
		shelfObjID, err := primitive.ObjectIDFromHex(shelfID)
		if err != nil {
			return "", fmt.Errorf("shelf_id is invalid")
		}

		objShelfID = append(objShelfID, shelfObjID)
	}

	location := &Location{
		ID:              primitive.NewObjectID(),
		LocationName:    req.LocationName,
		LocationAddress: req.LocationAddress,
		LocationImage:   req.LocationImage,
		LocationMap:     req.LocationMap,
		WareHouseID:     wareHouseObjID,
		BuildingID:      buildingObjID,
		FloorID:         floorObjID,
		RoomID:          roomObjID,
		ShelfID:         objShelfID,
		IsActice:        true,
		CreatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	locationID, err := s.repository.Create(ctx, location)
	if err != nil {
		return "", err
	}

	return locationID, nil

}

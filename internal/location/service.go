package location

import (
	"context"
	"fmt"
	"inventory-service/internal/product"
	"inventory-service/internal/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationService interface {
	CreateLocation(ctx context.Context, req *CreateLocationRequest, userID string) (string, error)
	GetLocations(ctx context.Context) ([]*LocationResponse, error)
	GetLocationByID(ctx context.Context, id string) (*LocationResponse, error)
	UpdateLocation(ctx context.Context, req *UpdateLocationRequest, id string) error
	DeleteLocation(ctx context.Context, id string) error
	AddProductToLocation(ctx context.Context, req *AddProductToLocationRequest, id string) error
}

type locationService struct {
	repository        LocationRepository
	storageRepository storage.StorageRepository
	productService    product.ProductService
}

func NewLocationService(repository LocationRepository, storageRepository storage.StorageRepository, productService product.ProductService) LocationService {
	return &locationService{
		repository:        repository,
		storageRepository: storageRepository,
		productService:    productService,
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

	floorObjID, err := primitive.ObjectIDFromHex(req.FloorID)
	if err != nil {
		return "", fmt.Errorf("floor_id is invalid")
	}

	roomObjID, err := primitive.ObjectIDFromHex(req.RoomID)
	if err != nil {
		return "", fmt.Errorf("room_id is invalid")
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
		Products:        []Product{},
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

func (s *locationService) GetLocations(ctx context.Context) ([]*LocationResponse, error) {
	
	locations, err := s.repository.GetLocations(ctx)
	if err != nil {
		return nil, err
	}

	var data []*LocationResponse
	for _, location := range locations {

		wareHouseData, err := s.storageRepository.GetStorageByID(ctx, &location.WareHouseID)
		if err != nil {
			return nil, err
		}
		warehouseInfo := &WarehouseInfo{
			ID:   wareHouseData.ID.Hex(),
			Name: wareHouseData.Name,
		}

		buildingData, err := s.storageRepository.GetStorageByID(ctx, &location.BuildingID)
		if err != nil {
			return nil, err
		}
		buildingInfo := &BuildingInfo{
			ID:   buildingData.ID.Hex(),
			Name: buildingData.Name,
		}

		floorData, err := s.storageRepository.GetStorageByID(ctx, &location.FloorID)
		if err != nil {
			return nil, err
		}
		floorInfo := &FloorInfo{
			ID:   floorData.ID.Hex(),
			Name: floorData.Name,
		}

		roomData, err := s.storageRepository.GetStorageByID(ctx, &location.RoomID)
		if err != nil {
			return nil, err
		}
		roomInfo := &RoomInfo{
			ID:   roomData.ID.Hex(),
			Name: roomData.Name,
		}

		var shelves []*ShelfInfo
		for _, shelfID := range location.ShelfID {
			shelfData, err := s.storageRepository.GetStorageByID(ctx, &shelfID)
			if err != nil {
				return nil, err
			}
			shelves = append(shelves, &ShelfInfo{
				ID:   shelfData.ID.Hex(),
				Name: shelfData.Name,
			})
		}

		var products []*ProductInfo
		for _, product := range location.Products {
			productData, err := s.productService.GetProductByID(ctx, product.ProductID.Hex())
			if err != nil {
				return nil, err
			}

			if productData == nil {
				products = append(products, &ProductInfo{
					ID:           "",
					Name:         "",
					PriceStore:   0,
					PriceService: 0,
					Description:  "",
					Image:        "",
					FolderName:   "",
					TopicName:    "",
					Quantity:     product.Quantity,
				})
				continue
			}

			products = append(products, &ProductInfo{
				ID:           productData.ID,
				Name:         productData.Name,
				PriceStore:   productData.PriceStore,
				PriceService: productData.PriceService,
				Description:  productData.Description,
				Image:        productData.Image,
				FolderName:   productData.FolderName,
				TopicName:    productData.TopicName,
				Quantity:     product.Quantity,
			})
		}

		locationResp := &LocationResponse{
			ID:              location.ID.Hex(),
			LocationName:    location.LocationName,
			LocationAddress: location.LocationAddress,
			LocationImage:   location.LocationImage,
			LocationMap:     location.LocationMap,
			Warehouse:       warehouseInfo,
			Building:        buildingInfo,
			Floor:           floorInfo,
			Room:            roomInfo,
			Shelves:         shelves,
			Products:        products,
			CreatedBy:       location.CreatedBy,
			CreatedAt:       location.CreatedAt,
			UpdatedAt:       location.UpdatedAt,
		}

		data = append(data, locationResp)
	}

	return data, nil
}

func (s *locationService) GetLocationByID(ctx context.Context, id string) (*LocationResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("location_id is invalid")
	}

	location, err := s.repository.GetLocationByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if location == nil {
		return nil, fmt.Errorf("location not found")
	}

	wareHouseData, err := s.storageRepository.GetStorageByID(ctx, &location.WareHouseID)
	if err != nil {
		return nil, err
	}

	warehouseInfo := &WarehouseInfo{
		ID:   wareHouseData.ID.Hex(),
		Name: wareHouseData.Name,
	}

	buildingData, err := s.storageRepository.GetStorageByID(ctx, &location.BuildingID)
	if err != nil {
		return nil, err
	}

	buildingInfo := &BuildingInfo{
		ID:   buildingData.ID.Hex(),
		Name: buildingData.Name,
	}

	floorData, err := s.storageRepository.GetStorageByID(ctx, &location.FloorID)
	if err != nil {
		return nil, err
	}

	floorInfo := &FloorInfo{
		ID:   floorData.ID.Hex(),
		Name: floorData.Name,
	}

	roomData, err := s.storageRepository.GetStorageByID(ctx, &location.RoomID)
	if err != nil {
		return nil, err
	}

	roomInfo := &RoomInfo{
		ID:   roomData.ID.Hex(),
		Name: roomData.Name,
	}

	var shelves []*ShelfInfo

	for _, shelfID := range location.ShelfID {
		shelfData, err := s.storageRepository.GetStorageByID(ctx, &shelfID)
		if err != nil {
			return nil, err
		}

		shelfInfo := &ShelfInfo{
			ID:   shelfData.ID.Hex(),
			Name: shelfData.Name,
		}

		shelves = append(shelves, shelfInfo)
	}

	var products []*ProductInfo

	for _, productID := range location.Products {

		productData, err := s.productService.GetProductByID(ctx, productID.ProductID.Hex())
		if err != nil {
			return nil, err
		}

		if productData == nil {
			products = append(products, &ProductInfo{
				ID:           "",
				Name:         "",
				PriceStore:   0,
				PriceService: 0,
				Description:  "",
				Image:        "",
				FolderName:   "",
				TopicName:    "",
				Quantity:     productID.Quantity,
			})
			continue
		}

		products = append(products, &ProductInfo{
			ID:           productData.ID,
			Name:         productData.Name,
			PriceStore:   productData.PriceStore,
			PriceService: productData.PriceService,
			Description:  productData.Description,
			Image:        productData.Image,
			FolderName:   productData.FolderName,
			TopicName:    productData.TopicName,
			Quantity:     productID.Quantity,
		})
	}

	locationData := &LocationResponse{
		ID:              location.ID.Hex(),
		LocationName:    location.LocationName,
		LocationAddress: location.LocationAddress,
		LocationImage:   location.LocationImage,
		LocationMap:     location.LocationMap,
		Warehouse:       warehouseInfo,
		Building:        buildingInfo,
		Floor:           floorInfo,
		Room:            roomInfo,
		Shelves:         shelves,
		Products:        products,
		CreatedBy:       location.CreatedBy,
		CreatedAt:       location.CreatedAt,
		UpdatedAt:       location.UpdatedAt,
	}

	return locationData, nil

}

func (s *locationService) UpdateLocation(ctx context.Context, req *UpdateLocationRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("location_id is invalid")
	}

	location, err := s.repository.GetLocationByID(ctx, objectID)
	if err != nil {
		return err
	}

	if location == nil {
		return fmt.Errorf("location not found")
	}

	if req.LocationName != "" {
		location.LocationName = req.LocationName
	}

	if req.LocationAddress != "" {
		location.LocationAddress = req.LocationAddress
	}

	if req.LocationImage != "" {
		location.LocationImage = req.LocationImage
	}

	if req.LocationMap != "" {
		location.LocationMap = req.LocationMap
	}

	if req.WareHouseID != "" {
		whID, err := primitive.ObjectIDFromHex(req.WareHouseID)
		if err != nil {
			return fmt.Errorf("warehouse_id is invalid")
		}
		location.WareHouseID = whID
	}

	if req.BuildingID != "" {
		bID, err := primitive.ObjectIDFromHex(req.BuildingID)
		if err != nil {
			return fmt.Errorf("building_id is invalid")
		}
		location.BuildingID = bID
	}

	if req.FloorID != "" {
		fID, err := primitive.ObjectIDFromHex(req.FloorID)
		if err != nil {
			return fmt.Errorf("floor_id is invalid")
		}
		location.FloorID = fID
	}

	if req.RoomID != "" {
		rID, err := primitive.ObjectIDFromHex(req.RoomID)
		if err != nil {
			return fmt.Errorf("room_id is invalid")
		}
		location.RoomID = rID
	}

	if req.ShelfID != nil {
		var shelfIDs []primitive.ObjectID
		for _, sID := range req.ShelfID {
			oid, err := primitive.ObjectIDFromHex(sID)
			if err != nil {
				return fmt.Errorf("invalid shelf_id: %s", sID)
			}
			shelfIDs = append(shelfIDs, oid)
		}
		location.ShelfID = shelfIDs
	}

	err = s.repository.UpdateLocation(ctx, location)
	if err != nil {
		return err
	}

	return nil
}

func (s *locationService) DeleteLocation(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("location_id is invalid")
	}

	return s.repository.DeleteLocation(ctx, objectID)

}

func (s *locationService) AddProductToLocation(ctx context.Context, req *AddProductToLocationRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("location_id is invalid")
	}

	location, err := s.repository.GetLocationByID(ctx, objectID)
	if err != nil {
		return err
	}

	if location == nil {
		return fmt.Errorf("location not found")
	}

	if req.ProductID == "" {
		return fmt.Errorf("product_id is required")
	}

	if req.Quantity == 0 {
		return fmt.Errorf("quantity is required")
	}

	return s.repository.AddProductToLocation(ctx, objectID, req.ProductID, req.Quantity)
}

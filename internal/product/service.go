package product

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-service/pkg/constants"
	"inventory-service/pkg/consul"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

type ProductService interface {
	GetProductByID(ctx context.Context, id string) (*Product, error)
}

type productService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	productServiceStr = "product-service"
)

func NewProductService(client *api.Client) ProductService {
	mainServiceAPI := NewServiceAPI(client, productServiceStr)
	return &productService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		fmt.Printf("Error creating service discovery: %v\n", err)
		return nil
	}

	var service *api.CatalogService

	for i := 0; i < 10; i++ {
		service, err = sd.DiscoverService()
		if err == nil && service != nil {
			break
		}
		fmt.Printf("Waiting for service %s... retry %d/10\n", serviceName, i+1)
		time.Sleep(3 * time.Second)
	}

	if service == nil {
		fmt.Printf("Service %s not found after retries, continuing anyway...\n", serviceName)
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		fmt.Println("Running in LOCAL_TEST mode — overriding service address to localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*Product, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	product, err := s.client.getProductByID(id, token)

	if err != nil {
		if sc, ok := product["status_code"].(float64); ok && int(sc) == 500 {
			return nil, nil
		}
		return nil, err
	}

	innerData, ok := product["data"].(map[string]interface{})
	if !ok || innerData == nil {
		return nil, nil
	}

	idVal, _ := innerData["id"].(string)
	if idVal == "" {
		return nil, nil
	}

	nameVal, _ := innerData["product_name"].(string)

	priceStore, _ := innerData["original_price_store"].(float64)
	priceService, _ := innerData["original_price_service"].(float64)
	image, _ := innerData["cover_image"].(string)
	product_description, _ := innerData["product_description"].(string)

	topicData, _ := innerData["topic"].(map[string]interface{})
	folderData, _ := innerData["folder"].(map[string]interface{})

	topicName := ""
	folderName := ""

	if topicData != nil {
		topicName, _ = topicData["name"].(string)
	}

	if folderData != nil {
		folderName, _ = folderData["name"].(string)
	}

	return &Product{
		ID:           idVal,
		Name:         nameVal,
		PriceStore:   priceStore,
		PriceService: priceService,
		Description:  product_description,
		Image:        image,
		FolderName:   folderName,
		TopicName:    topicName,
	}, nil

}

func (c *callAPI) getProductByID(id string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/api/v1/products/%s", id)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var productData interface{}

	err = json.Unmarshal([]byte(res), &productData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := productData.(map[string]interface{})

	return myMap, nil
}

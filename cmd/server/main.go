package main

import (
	"context"
	"fmt"
	"inventory-service/config"

	// "inventory-service/internal/product"
	productplacement "inventory-service/internal/product_placement"
	producttransaction "inventory-service/internal/product_transaction"
	shelfquantity "inventory-service/internal/shelf_quantity"
	shelftype "inventory-service/internal/shelf_type"
	"inventory-service/internal/storage"
	"inventory-service/pkg/consul"
	"inventory-service/pkg/uploader"
	"inventory-service/pkg/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load env
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Successfully loaded .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	if err := waitPassing(consulClient, "product-service", 60*time.Second); err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Handle OS signal để deregister
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server... De-registering from Consul...")
		consulConn.Deregister()
		os.Exit(0)
	}()

	imageService := uploader.NewImageService(consulClient)

	shelfTypeCollection := mongoClient.Database(cfg.MongoDB).Collection("shelf_type")
	storageCollection := mongoClient.Database(cfg.MongoDB).Collection("storage")
	productTransaction := mongoClient.Database(cfg.MongoDB).Collection("product_transaction")
	productPlacement := mongoClient.Database(cfg.MongoDB).Collection("product_placement")
	shelfQuantityCollection := mongoClient.Database(cfg.MongoDB).Collection("shelf_quantity")
	shelfTypeRepository := shelftype.NewShelfTypeRepository(shelfTypeCollection)
	storageRepository := storage.NewStorageRepository(storageCollection)
	productTransactionRepository := producttransaction.NewProductTransactionRepository(productTransaction)
	productPlacementRepository := productplacement.NewProductPlacementRepository(productPlacement)

	shelfQuantityRepository := shelfquantity.NewShelfQuantityRepository(shelfQuantityCollection)
	shelfQuantityService := shelfquantity.NewShelfQuantityService(shelfQuantityRepository, storageRepository)
	shelfQuantityHandler := shelfquantity.NewShelfQuantityHandler(shelfQuantityService)

	shelfTypeService := shelftype.NewShelfTypeService(shelfTypeRepository, storageRepository, imageService)
	shelfTypeHandler := shelftype.NewShelfTypeHandler(shelfTypeService)

	// productService := product.NewProductService(consulClient)
	storageService := storage.NewStorageService(storageRepository, shelfTypeRepository, shelfQuantityRepository, imageService)
	storageHandler := storage.NewStorageHandler(storageService)

	productPlacementService := productplacement.NewProductPlacementService(productPlacementRepository, storageRepository)
	productPlacementHandler := productplacement.NewProductPlacementHandler(productPlacementService)
	productTransactionService := producttransaction.NewProductTransactionService(productTransactionRepository, productPlacementService, mongoClient)
	productTransactionHandler := producttransaction.NewProductTransactionHandler(productTransactionService)

	r := gin.Default()
	shelftype.RegisterRoutes(r, shelfTypeHandler)
	storage.RegisterRoutes(r, storageHandler)
	productplacement.RegisterRoutes(r, productPlacementHandler)
	producttransaction.RegisterRoutes(r, productTransactionHandler)
	shelfquantity.RegisterRoutes(r, shelfQuantityHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8009"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func waitPassing(cli *consulapi.Client, name string, timeout time.Duration) error {
	dl := time.Now().Add(timeout)
	for time.Now().Before(dl) {
		entries, _, err := cli.Health().Service(name, "", true, nil)
		if err == nil && len(entries) > 0 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("%s not ready in consul", name)
}

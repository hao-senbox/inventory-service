package main

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/location"
	"inventory-service/internal/product"
	"inventory-service/internal/storage"
	"inventory-service/pkg/consul"
	"inventory-service/pkg/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()


	// Handle OS signal để deregister
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server... De-registering from Consul...")
		consulConn.Deregister()
		os.Exit(0)
	}()
	productService := product.NewProductService(consulClient)
	storageCollection := mongoClient.Database(cfg.MongoDB).Collection("storage")
	storageRepository := storage.NewStorageRepository(storageCollection)
	storageService := storage.NewStorageService(storageRepository)
	storageHandler := storage.NewStorageHandler(storageService)

	locationCollection := mongoClient.Database(cfg.MongoDB).Collection("location")
	locationRepository := location.NewLocationRepository(locationCollection)
	locationService := location.NewLocationService(locationRepository, storageRepository, productService)
	locationHandler := location.NewLocationHandler(locationService)

	r := gin.Default()

	storage.RegisterRoutes(r, storageHandler)
	location.RegisterRoutes(r, locationHandler)
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

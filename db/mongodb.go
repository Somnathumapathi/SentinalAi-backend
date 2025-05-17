package db

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var (
// 	Client        *mongo.Client
// 	Database      *mongo.Database
// 	Subscriptions *mongo.Collection
// 	Clients       *mongo.Collection
// )

// func InitDB() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Get MongoDB URI from environment variable
// 	mongoURI := os.Getenv("MONGODB_URI")
// 	if mongoURI == "" {
// 		mongoURI = "mongodb://localhost:27017"
// 	}

// 	// Connect to MongoDB
// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to MongoDB: %v", err)
// 	}

// 	// Ping the database
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to ping MongoDB: %v", err)
// 	}

// 	// Set up database and collections
// 	Client = client
// 	Database = client.Database("sentinal")
// 	Subscriptions = Database.Collection("subscriptions")
// 	Clients = Database.Collection("clients")

// 	log.Println("Connected to MongoDB!")
// 	return nil
// }

// func CloseDB() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := Client.Disconnect(ctx); err != nil {
// 		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
// 	}

// 	return nil
// }

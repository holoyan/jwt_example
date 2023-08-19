package core

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func Connect() *mongo.Client {

	host := os.Getenv("MONGO_HOST")
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")

	// Set up MongoDB client options with authentication
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017/?authSource=admin", username, password, host))

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client

	// Close the connection when done
	//defer client.Disconnect(context.Background())
}

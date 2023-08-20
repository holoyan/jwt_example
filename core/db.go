package core

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

type FindOneResult struct {
	Success bool
	Message string
	Result map[string]interface{}
}

func FindOne(coll string, key string, value string) FindOneResult {
	var result map[string]interface{}

	var response FindOneResult

	response.Success = true
	filter := bson.M{key: value}
	if key == "_id" {
		objectId, _ := primitive.ObjectIDFromHex(value)
		filter[key] = objectId
	}

	err := DB.Collection(coll).FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		response.Success = false
		if err == mongo.ErrNoDocuments {
			response.Message = "Not found"
			return response
		}
		response.Message = "Something went wrong"
		return response
	}

	response.Result = result

	return response
}

func DeleteOne(coll string, key string, value string) bool {
	filter := bson.M{key: value}
	if key == "_id" {
		objectId, _ := primitive.ObjectIDFromHex(value)
		filter[key] = objectId
	}

	deleteResult, err := DB.Collection(coll).DeleteOne(context.Background(), filter)

	if err != nil {
		return false
	}

	return deleteResult.DeletedCount > 0
}

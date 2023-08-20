package core

import (
	"context"
	"fmt"
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

//func findOne(coll string, key string, value string) {
//	var result map[string]interface{}
//	objectID, err1 := primitive.ObjectIDFromHex(value)
//	filter := bson.M{"_id": objectID}
//	if err1 != nil {
//		log.Fatal(err1)
//	}
//
//	err := DB.Collection(coll).FindOne(context.Background(), filter).Decode(&result)
//
//	if err != nil {
//		// ErrNoDocuments means that the filter did not match any documents in the collection
//		response["success"] = false
//		if err == mongo.ErrNoDocuments {
//			response["message"] = "Not found"
//			core.JsonResponse(res, response, 404)
//			return
//		}
//		response["message"] = "Something went wrong"
//		core.JsonResponse(res, response, 400)
//		return
//	}
//}

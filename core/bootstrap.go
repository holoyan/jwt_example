package core

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var DbClient *mongo.Client
var DB *mongo.Database

func Load()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DbClient = Connect()
	DB = DbClient.Database(os.Getenv("MONGO_DATABASE"))
}

func Close()  {
	defer DbClient.Disconnect(context.Background())
}

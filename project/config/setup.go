package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"golangdemo/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB Client instance
var DB *mongo.Client

// ConnectDB function is used to instantiate MongoDB Connection
func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoHost := os.Getenv("MONGO_HOST")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := url.QueryEscape(os.Getenv("MONGO_PASSWORD"))
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, mongoHost)
	fmt.Println(mongoURI)
	//encodedMongoUri := url.QueryEscape(mongoURI)
	//completeURI := fmt.Sprintf("mongodb://%s/Initializ?retryWrites=true&w=majority", encodedMongoUri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	DB = client
}

func GetRepoCollection(collectionName string) repository.Repository {
	repo := repository.MongoUserRepository{
		Collection: GetCollection(DB, collectionName),
	}
	return &repo
}

// GetCollection function helps in getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Initializ-dev").Collection(collectionName)
	return collection
}

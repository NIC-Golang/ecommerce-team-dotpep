package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Print("Error with connecting to Mongo!")
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Print("Cannot ping MongoDB!")
		log.Fatal(err)
	}

	log.Print("Successfully connected to MongoDB")
	return client
}

var DB *mongo.Client = ConnectMongo()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Cluster0").Collection(collectionName)
	return collection
}

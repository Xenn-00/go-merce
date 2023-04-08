package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSetup() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(""))
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection

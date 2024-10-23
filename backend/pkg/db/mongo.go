package db

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// ConnectDB accepts a MongoDB URI as a parameter
func ConnectDB(mongoURI string) {
    clientOptions := options.Client().ApplyURI(mongoURI)
    var err error
    mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = mongoClient.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB!")
}

// GetCollection retrieves the specified collection from the MongoDB database
func GetCollection(name string) *mongo.Collection {
    return mongoClient.Database("dockerTracker").Collection(name)
}

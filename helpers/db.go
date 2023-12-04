package helpers

import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Usercollection *mongo.Collection
var client *mongo.Client

func ConnectToDB() {
	connectionurl := MONGO_URL
	clientOptions := options.Client().ApplyURI(connectionurl)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create a MongoDB client
	client, ERR := mongo.Connect(ctx, clientOptions)
	if ERR != nil {
		log.Fatal(ERR)
	}

	// Ping the MongoDB server to check the connection
	ERR = client.Ping(ctx, nil)
	if ERR != nil {
		log.Fatal(ERR)
	}

	fmt.Println("Connected to MongoDB!")

	Usercollection = client.Database("Feedback").Collection("user")
}

func DisconnectFromDB() {
	
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Disconnected from MongoDB Atlas!")
}

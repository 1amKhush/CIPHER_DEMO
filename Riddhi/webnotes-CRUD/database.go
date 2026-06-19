// Connects to MongoDB and sets a Stable API version
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client // Global variable to hold the MongoDB client

func ConnectDB() {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")


	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pinged MongoDB!")

	
}	





package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var NoteCollection *mongo.Collection

func Connect_db() {
	uri := "mongodb+srv://campy:campy123@cluster0.xuumolv.mongodb.net/?appName=Cluster0"
	if uri == "" {
		log.Fatal("MONGOURI couldnt be accessed")
	}

	clientOptions := options.Client().
		ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	fmt.Println("DB connection successful")

	UserCollection = client.Database("notehub").Collection("note_users")
	NoteCollection = client.Database("notehub").Collection("notes")
}

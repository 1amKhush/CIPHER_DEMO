package controller

import (
	"CRUD-project/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var userCollection *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	connectionString := os.Getenv("URI")
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection successful")
	collection = client.Database("notesDB").Collection("notes")
	userCollection = client.Database("notesDB").Collection("users")
}
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Failed")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Failed")
			return
		}
		token := parts[1]
		if token != "1818" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Failed")
			return
		}
		next(w, r)
	}
}
func CreateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note model.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		json.NewEncoder(w).Encode("Failed")
		return
	}
	_, err = collection.InsertOne(context.Background(), note)
	if err != nil {
		json.NewEncoder(w).Encode("Failed")
		return
	}
	json.NewEncoder(w).Encode("Completed")
}
func ReadNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var allNotes []model.Note
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode("Failed")
		return
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.Background(), &allNotes)
	if err != nil {
		json.NewEncoder(w).Encode("Failed")
		return
	}
	json.NewEncoder(w).Encode(allNotes)
}
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedNote model.Note
	err := json.NewDecoder(r.Body).Decode(&updatedNote)
	if err != nil {
		json.NewEncoder(w).Encode("Error")
		return
	}
	filter := bson.M{"_id": updatedNote.ID}
	update := bson.M{
		"$set": bson.M{
			"content": updatedNote.Content,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		json.NewEncoder(w).Encode("Error")
		return
	}
	json.NewEncoder(w).Encode("Completed")
}
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note model.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		json.NewEncoder(w).Encode("Error")
		return
	}
	filter := bson.M{"_id": note.ID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		json.NewEncoder(w).Encode("Error")
		return
	}
	json.NewEncoder(w).Encode("Completed")
}

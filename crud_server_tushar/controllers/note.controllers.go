package controllers

import (
	"context"
	"crud_server/database"
	"crud_server/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateNoteReq struct {
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedBy primitive.ObjectID `json:"createdBy"`
}
type UpdateNoteReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func Homeroute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello goluuu<h1/>"))
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note

	params := mux.Vars(r)
	userId := params["id"]
	newUserId, _ := primitive.ObjectIDFromHex(userId)

	fmt.Println(newUserId)

	err := database.NoteCollection.FindOne(context.Background(), bson.M{
		"_id": newUserId,
	}).Decode(&note)

	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	fmt.Printf("The note having this id : %+v is : ", note)

}

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	cur, err := database.NoteCollection.Find(context.Background(), bson.M{}) // you can write bson.M{} or bson.D{{}}, becoz both of them have their own structure
	// bson.M{
	//     "email": "abc@gmail.com",
	// }
	// 	bson.D{
	//     {"email", "abc@gmail.com"},
	// }
	// bson.M is a map, whereas bson.D is like an ordered list of objects, and now u can understand why each of these work

	if err != nil {
		http.Error(w, "not a single note exists", http.StatusNotFound)
	}

	// var notes []bson.D
	var notes []models.Note

	// 1. This is the longer method where we decode each element one at a time

	for cur.Next(context.Background()) {
		// var note bson.D
		var note models.Note
		err := cur.Decode(&note)

		if err != nil {
			http.Error(w, "idk what error this is", http.StatusNotFound)
		}

		notes = append(notes, note)

		fmt.Printf("\n%#v", note)
	}

	// 2. This is the shorter way to decode all the elements from the cursor

	// cur.All(context.Background(), &notes)

	// fmt.Print(notes)

	defer cur.Close(context.Background())
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("user").(*Claims)

	// fmt.Println("\n\n", claims)

	userCreatingNoteId := claims.RegisteredClaims.Subject

	req.CreatedBy, err = primitive.ObjectIDFromHex(userCreatingNoteId)

	database.NoteCollection.InsertOne(context.Background(), req)

	fmt.Println("New Note created :", req)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	var req UpdateNoteReq

	params := mux.Vars(r)
	noteId := params["id"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	newnoteId, _ := primitive.ObjectIDFromHex(noteId)

	claims := r.Context().Value("user").(*Claims)

	userUpdatingNoteId := claims.RegisteredClaims.Subject

	curuserUpdatingNoteId, err := primitive.ObjectIDFromHex(userUpdatingNoteId)

	// fmt.Println(curuserUpdatingNoteId)
	// fmt.Println(newnoteId)

	filter := bson.M{
		"_id":       newnoteId,
		"createdby": curuserUpdatingNoteId,
	}

	update := bson.M{
		"$set": bson.M{
			"title":   req.Title,
			"content": req.Content,
		},
	}

	result, err := database.NoteCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		http.Error(
			w,
			"update failed",
			http.StatusInternalServerError,
		)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(
			w,
			"note not found",
			http.StatusNotFound,
		)
		return
	}

	fmt.Println("Updated Note : ", result)

}

func DeleteNote(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userId := params["id"]
	newUserId, _ := primitive.ObjectIDFromHex(userId)

	filter := bson.M{
		"_id": newUserId,
	}

	result, err := database.NoteCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	fmt.Println("The note deleted is : ", result)
}

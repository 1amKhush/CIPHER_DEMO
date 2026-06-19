package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func getJWTSecret() []byte {
    secret := os.Getenv("jwt")
    if secret == "" {
        log.Fatal("JWT secret is not set")
    }
    return []byte(secret)
}

func GenerateJWT(email string) (string, error) {
	claims:=jwt.MapClaims{ // claim means the data we want to include in the token
		"email": email,
		"exp":jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // token expires in 24 hours
		"iat":jwt.NewNumericDate(time.Now()), // issued at time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func Authmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
		},
		)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	
		next.ServeHTTP(w, r)
	})
}

func user(w http.ResponseWriter , r *http.Request) {
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	var user User
	err:= json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	
	hash, err := HashPassword(user.Password)
    if err != nil {
    log.Fatal(err)
     }

    user.Password = hash

    createuser, err := collection.InsertOne(context.TODO(), user)
    if err != nil {
    log.Fatal(err)
    }

	fmt.Println("Inserted a single document: ", createuser.InsertedID)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
	
}

func login(w http.ResponseWriter , r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	
	var login Login
	err:= json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
	w.WriteHeader(http.StatusUnauthorized)
	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}
	
	var user User
	err = collection.FindOne(context.TODO(), bson.M{"email": login.Email}).Decode(&user)
	if err != nil {
	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}

	if CheckPasswordHash(login.Password, user.Password) {
    token, err := GenerateJWT(user.Email)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	    w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Login successful",
			"token":   token,
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid email or password",
		})
	}
}

func createNote(note Note) {
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	createResult, err := collection.InsertOne(context.TODO(), note)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", createResult.InsertedID)
}

func updateNote(id string,note Note) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	filter := bson.M{"_id": objID}
	
	update := bson.M{"$set": bson.M{"title": note.Title, "content": note.Content}}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func deleteNote(id string) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	filter := bson.M{"_id": objID}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the notes collection\n", deleteResult.DeletedCount)
}

func deleteAllNotes() int64{
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the notes collection\n", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

func getNote(id string) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	var result Note
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
}

func getAllNotes() []Note{
	collection := client.Database(os.Getenv("dbname")).Collection(os.Getenv("collectionname"))
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var results []Note 
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found multiple documents: %+v\n", results)
	return results
}

// actual controller

func GetAllNotes(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type","application/json")
	allNotes:=getAllNotes()
	json.NewEncoder(w).Encode(allNotes)
}

func CreateNote(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	var note Note
	_ = json.NewDecoder(r.Body).Decode(&note)
	createNote(note)
	json.NewEncoder(w).Encode(note)
}

func UpdateNote(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	params:=mux.Vars(r)
	var note Note
    json.NewDecoder(r.Body).Decode(&note)
	updateNote(params["id"], note)
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteNote(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	params:=mux.Vars(r)
	deleteNote(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteALLNotes(w http.ResponseWriter , r *http.Request){
	{
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	count:=deleteAllNotes()
	json.NewEncoder(w).Encode(count)
}
}


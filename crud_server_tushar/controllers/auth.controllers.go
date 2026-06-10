package controllers

import (
	"context"
	"crud_server/database"
	"crud_server/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	// ID    string `json:"_id"`    // this is a useless field bcoz while issuing tokens, the claims that you created
	// had the userid in the subject of registered claims, and the email, thats it, so theres no point of creating ID here

}

type TokenJsonResponse struct {
	AccessToken string `json:"access_token"`
}

var JwtSecret = []byte("hehe")

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func issueAccessToken(user *models.User) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(45 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: user.Email,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtSecret)
}

func issueRefreshToken(user *models.User) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JwtSecret)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&req)

	fmt.Println(req)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "email, password and username are required", http.StatusBadRequest)
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": req.Username},
		},
	}

	err = database.UserCollection.FindOne(context.Background(), filter).Decode(&user)

	if err == nil {
		http.Error(w, "username/email already exists", http.StatusConflict)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	if err != nil {
		log.Fatal("Some problem in hasing the passwd")
	}

	req.Password = string(hashedPass)

	fmt.Println(req)

	inserted, err := database.UserCollection.InsertOne(context.Background(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Registered a new user with id : ", inserted.InsertedID)

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&req)

	fmt.Println(req)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate fields explicitly
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	// Email format verification
	// if !isValidEmail(req.Email) {
	// 	http.Error(w, "invalid email format", http.StatusBadRequest)
	// 	return
	// }

	filter := bson.M{
		"email": req.Email,
	}

	database.UserCollection.FindOne(context.Background(), filter).Decode(&user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := issueAccessToken(&user)
	if err != nil {
		http.Error(w, "internal error in access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := issueRefreshToken(&user)
	if err != nil {
		http.Error(w, "internal error in refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // avoids CSRF attacks, agar kisi aur site ke same endpt ne hit kiya toh usko bhi jaa sakta hai cookie
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/auth/refresh", // yeh check krlena bcoz this cookie containing refresh token will be sent
		// only to this path
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenJsonResponse{AccessToken: accessToken})

	fmt.Println("Logged in a user with id :", user.ID)

}

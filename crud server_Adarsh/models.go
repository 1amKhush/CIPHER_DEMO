package main

import "github.com/golang-jwt/jwt/v5"


type Task struct 
{ 
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
}

type User struct 
{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"pwd"`
}


type JWTClaims struct {

	Email string `json:"email"`
	Name string `json:"name"`
	jwt.RegisteredClaims

	// Add more fields as needed
}


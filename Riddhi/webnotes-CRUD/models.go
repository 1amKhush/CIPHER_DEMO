package main

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title string `bson:"title,omitempty" json:"title,omitempty"`
	Content string `bson:"content,omitempty" json:"content,omitempty"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Email string `bson:"email,omitempty" json:"email,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
}

type Login struct {
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
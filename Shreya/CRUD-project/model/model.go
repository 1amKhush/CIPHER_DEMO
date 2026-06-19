package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"-" bson:"password,omitempty"`
}

type Note struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
}

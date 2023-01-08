package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Recipe struct {
	// ID           string    `json:"id"`
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"publishedAt" bson:"publishedAt"`
}
type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
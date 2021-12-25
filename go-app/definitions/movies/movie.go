package movies

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Movie struct
type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
	Cover       string             `bson:"cover,omitempty"`
	AddedBy     primitive.ObjectID `bson:"added_by,omitempty"`
}

type AddMovieInput struct {
	Name        string    `json:"name" mod:"trim,title" binding:"required"`
	Description string    `json:"description" mod:"trim" binding:"required"`
	Date        time.Time `json:"date"`
}
type AddMovieOutput struct {
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
	Cover       string             `bson:"cover"`
	AddedBy     primitive.ObjectID `bson:"added_by"`
}

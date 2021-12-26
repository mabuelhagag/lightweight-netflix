package movies

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"
)

// Movie struct
type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Date        time.Time          `bson:"date,omitempty"`
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
type UploadCoverInput struct {
	Cover *multipart.FileHeader `form:"cover" binding:"required"`
}

type UpdateMovieInput struct {
	Name        string    `json:"name" mod:"trim,title"`
	Description string    `json:"description" mod:"trim"`
	Date        time.Time `json:"date"`
}

type WatchedMovieEntry struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	MovieID primitive.ObjectID `bson:"movie_id"`
	UserId  primitive.ObjectID `bson:"user_id"`
	Time    time.Time          `bson:"time"`
}

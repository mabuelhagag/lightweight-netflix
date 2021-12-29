package movies

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"
)

// Movie struct
type Movie struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string             `bson:"name,omitempty"`
	Description      string             `bson:"description,omitempty"`
	Date             time.Time          `bson:"date,omitempty"` // TODO: use string to parse date from it
	AddedBy          primitive.ObjectID `bson:"added_by,omitempty"`
}

type AddMovieInput struct {
	Name        string    `json:"name" mod:"trim,title" binding:"required"`
	Description string    `json:"description" mod:"trim" binding:"required"`
	Date        time.Time `json:"date"` // TODO: use string to parse date from it
}
type AddMovieOutput struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"` // TODO: use string to parse date to it
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
	mgm.DefaultModel `bson:",inline"`
	MovieID          primitive.ObjectID `bson:"movie_id"`
	UserId           primitive.ObjectID `bson:"user_id"`
}

func (m *WatchedMovieEntry) CollectionName() string {
	return "watched"
}

type ReviewMovieInput struct {
	Rating uint8  `json:"rating" binding:"required,gte=1,lte=5"`
	Review string `json:"review" mod:"trim"`
}

type ReviewMovieEntry struct {
	mgm.DefaultModel `bson:",inline"`
	MovieID          primitive.ObjectID `bson:"movie_id"`
	UserId           primitive.ObjectID `bson:"user_id"`
	Rating           uint8              `bson:"rating"`
	Review           string             `bson:"review"`
}

func (m *ReviewMovieEntry) CollectionName() string {
	return "reviews"
}

type MovieInfo struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string    `bson:"name,omitempty"`
	Description      string    `bson:"description,omitempty"`
	Date             time.Time `bson:"date,omitempty"` // TODO: use string to parse date from it
	Rating           float32   `bson:"rating,omitempty"`
}

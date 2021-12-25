package moviesrepo

import (
	"context"
	"errors"
	"go-app/definitions/movies"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Repo Interface
type Repo interface {
	AddMovie(movie *movies.Movie) (*movies.Movie, error)
}

type moviesRepo struct {
	db *mongo.Client
}

// NewMoviesRepo will instantiate Movies Repository
func NewMoviesRepo(db *mongo.Client) Repo {
	return &moviesRepo{
		db: db,
	}
}

func (b *moviesRepo) AddMovie(movie *movies.Movie) (*movies.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("movies")

	result, err := collection.InsertOne(ctx, *movie)

	if err := collection.FindOne(ctx, bson.D{{"_id", result.InsertedID}}).Decode(&movie); err != nil {
		return nil, errors.New("Unable to get movie")
	}
	return movie, err
}

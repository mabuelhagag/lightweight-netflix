package moviesrepo

import (
	"context"
	"errors"
	"github.com/kamva/mgm/v3"
	"go-app/definitions/movies"
	"go-app/definitions/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Repo Interface
type Repo interface {
	AddMovie(movie *movies.Movie) (*movies.Movie, error)
	GetMovieById(id string) (*movies.Movie, error)
	UpdateMovie(movie *movies.Movie, id string) (*movies.Movie, error)
	DeleteMovie(id string) error
	AddToWatchedList(watchEntry *movies.WatchedMovieEntry) error
	DidWatchMovie(movie *movies.Movie, user *users.User) (bool, error)
	ReviewMovie(reviewEntry *movies.ReviewMovieEntry) error
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
	if err := mgm.Coll(movie).Create(movie); err != nil {
		return nil, err
	}

	return movie, nil
}

func (b *moviesRepo) GetMovieById(id string) (*movies.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("movies")

	var movie *movies.Movie
	movieID, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": movieID}}
	err = collection.FindOne(ctx, filter).Decode(&movie)
	if err != nil {
		return nil, errors.New("Unable to get movie")
	}
	return movie, err
}

func (b *moviesRepo) UpdateMovie(movie *movies.Movie, id string) (*movies.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("movies")

	objectId, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objectId}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.D{{"$set", movie}}
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Unable to get movie")
		}
		return nil, err
	}
	return movie, nil
}

func (b *moviesRepo) DeleteMovie(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // releases resources if CreateUser completes before timeout elapses
	collection := b.db.Database("lw-netflix").Collection("movies")

	objectId, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objectId}}
	var movie *movies.Movie
	err = collection.FindOneAndDelete(ctx, filter).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("Unable to get movie")
		}
		return err
	}
	return nil
}

func (b *moviesRepo) AddToWatchedList(watchEntry *movies.WatchedMovieEntry) error {

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", watchEntry.UserId}}}},
				bson.D{{"movie_id", bson.D{{"$eq", watchEntry.MovieID}}}},
			}},
	}
	err := mgm.Coll(watchEntry).First(filter, watchEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err := mgm.Coll(watchEntry).Create(watchEntry)
			if err != nil {
				return err
			}
		}
	}
	err = mgm.Coll(watchEntry).Update(watchEntry)
	if err != nil {
		return err
	}
	return nil
}

func (b *moviesRepo) DidWatchMovie(movie *movies.Movie, user *users.User) (bool, error) {
	watchEntry := &movies.WatchedMovieEntry{}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", user.ID}}}},
				bson.D{{"movie_id", bson.D{{"$eq", movie.ID}}}},
			}},
	}
	err := mgm.Coll(watchEntry).First(filter, watchEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (b *moviesRepo) ReviewMovie(reviewEntry *movies.ReviewMovieEntry) error {

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user_id", bson.D{{"$eq", reviewEntry.UserId}}}},
				bson.D{{"movie_id", bson.D{{"$eq", reviewEntry.MovieID}}}},
			}},
	}
	foundEntry := &movies.ReviewMovieEntry{}
	err := mgm.Coll(foundEntry).First(filter, foundEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err := mgm.Coll(reviewEntry).Create(reviewEntry)
			if err != nil {
				return err
			}
		}
	}
	reviewEntry.ID = foundEntry.ID
	err = mgm.Coll(reviewEntry).Update(reviewEntry)
	if err != nil {
		return err
	}
	return nil

}

package moviesrepo

import (
	"github.com/kamva/mgm/v3"
	"go-app/definitions/movies"
	"go-app/definitions/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repo Interface
type Repo interface {
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

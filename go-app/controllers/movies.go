package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/builder"
	"github.com/kamva/mgm/v3/operator"
	"go-app/definitions/movies"
	"go-app/definitions/users"
	"go-app/repositories/moviesrepo"
	"go-app/repositories/usersrepo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

// MoviesController interface
type MoviesController interface {
	AddMovie(*gin.Context)
	UploadCover(*gin.Context)
	UpdateMovie(*gin.Context)
	DeleteMovie(c *gin.Context)
	WatchMovie(c *gin.Context)
	ReviewMovie(c *gin.Context)
	ListWatchedMovies(c *gin.Context)
	ListMovies(c *gin.Context)
	GetMovieInfo(c *gin.Context)
}

type moviesController struct {
	mr moviesrepo.Repo
	ur usersrepo.Repo
}

// NewMoviesController instantiates User Controller
func NewMoviesController(br moviesrepo.Repo, us usersrepo.Repo) MoviesController {
	return &moviesController{mr: br, ur: us}
}

func (ctl *moviesController) AddMovie(c *gin.Context) {
	var movieInput movies.AddMovieInput
	if err := c.ShouldBindJSON(&movieInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	movie, err := ctl.inputToMovie(movieInput, c)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
	}
	if err := mgm.Coll(movie).Create(movie); err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Failed while adding movie", err.Error())
		return
	}

	output := ctl.movieToOutput(movie)
	HTTPRes(c, http.StatusOK, "Movie added", output)
}

func (ctl *moviesController) inputToMovie(input movies.AddMovieInput, c *gin.Context) (*movies.Movie, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	currentUser := c.MustGet("user").(*users.User)
	return &movies.Movie{
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
		AddedBy:     currentUser.ID,
	}, nil
}
func (ctl *moviesController) movieToOutput(movie *movies.Movie) *movies.AddMovieOutput {
	return &movies.AddMovieOutput{
		ID:          movie.ID.Hex(),
		Name:        movie.Name,
		Description: movie.Description,
		Date:        movie.Date,
	}
}

func (ctl *moviesController) UploadCover(c *gin.Context) {
	var uploadCoverInput movies.UploadCoverInput
	if err := c.ShouldBind(&uploadCoverInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", "Movie ID not provided")
		return
	}
	movie := &movies.Movie{}
	err := mgm.Coll(movie).FindByID(movieId, movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	err = c.SaveUploadedFile(uploadCoverInput.Cover, "/opt/go-app/covers/"+movieId+".jpg")
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "File upload error", err.Error())
		return
	}
	HTTPRes(c, http.StatusOK, "Cover uploaded", nil)
}

func (ctl *moviesController) UpdateMovie(c *gin.Context) {
	var movieInput movies.UpdateMovieInput
	if err := c.ShouldBindJSON(&movieInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", "Movie ID not provided")
		return
	}

	currentUser := c.MustGet("user").(*users.User)
	movie := &movies.Movie{}
	err := mgm.Coll(movie).FindByID(movieId, movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	if movie.AddedBy != currentUser.ID {
		HTTPRes(c, http.StatusForbidden, "Error updating movie info", "Movie is not owned by current user")
		return
	}
	err = ctl.updateMovieInputToMovie(movieInput, movie)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	err = mgm.Coll(movie).Update(movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Failed while updating movie", err.Error())
		return
	}

	_ = mgm.Coll(movie).FindByID(movieId, movie)
	output := ctl.movieToOutput(movie)
	HTTPRes(c, http.StatusOK, "Movie Updated", output)
}

func (ctl *moviesController) updateMovieInputToMovie(input movies.UpdateMovieInput, output *movies.Movie) error {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return err
	}
	output.Name = input.Name
	output.Description = input.Description
	output.Date = input.Date

	return nil
}

func (ctl *moviesController) DeleteMovie(c *gin.Context) {
	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}

	currentUser := c.MustGet("user").(*users.User)
	movie := &movies.Movie{}
	err := mgm.Coll(movie).FindByID(movieId, movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	if movie.AddedBy != currentUser.ID {
		HTTPRes(c, http.StatusForbidden, "Error updating movie info", "Movie is not owned by current user")
		return
	}
	err = mgm.Coll(movie).Delete(movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error deleting movie", err.Error())
		return
	}
	HTTPRes(c, http.StatusOK, "Movie Deleted", nil)
}
func (ctl *moviesController) WatchMovie(c *gin.Context) {
	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}

	currentUser := c.MustGet("user").(*users.User)
	movie := &movies.Movie{}
	err := mgm.Coll(movie).FindByID(movieId, movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	watchedEntry := movies.WatchedMovieEntry{MovieID: movie.ID, UserId: currentUser.ID}
	if err = ctl.mr.AddToWatchedList(&watchedEntry); err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error adding movie to watch list", err.Error())
		return
	}
	HTTPRes(c, http.StatusOK, "Movie added to watch list", nil)
}

func (ctl *moviesController) ReviewMovie(c *gin.Context) {
	var reviewInput movies.ReviewMovieInput
	if err := c.ShouldBindJSON(&reviewInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", "Movie ID not provided")
		return
	}

	currentUser := c.MustGet("user").(*users.User)

	movie := &movies.Movie{}
	err := mgm.Coll(movie).FindByID(movieId, movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			HTTPRes(c, http.StatusNotFound, "Movie not found", nil)
			return
		}
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}

	watchedMovie, err := ctl.mr.DidWatchMovie(movie, currentUser)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error while checking watched movie", err.Error())
		return
	}
	if watchedMovie != true {
		HTTPRes(c, http.StatusForbidden, "Error reviewing movie", "User hasn't watched the movie yet")
		return
	}

	reviewEntry, err := ctl.reviewMovieInputToReviewMovieEntry(reviewInput, currentUser, movie)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	err = ctl.mr.ReviewMovie(reviewEntry)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error while reviewing movie", err.Error())
		return
	}

	HTTPRes(c, http.StatusOK, "Movie reviewed", nil)

}

func (ctl *moviesController) reviewMovieInputToReviewMovieEntry(input movies.ReviewMovieInput, user *users.User, movie *movies.Movie) (*movies.ReviewMovieEntry, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	return &movies.ReviewMovieEntry{
		MovieID: movie.ID,
		UserId:  user.ID,
		Rating:  input.Rating,
		Review:  input.Review,
	}, nil
}

func (ctl *moviesController) ListWatchedMovies(c *gin.Context) {

	currentUser := c.MustGet("user").(*users.User)

	// TODO: aggregate to get movie details
	watchedMovies := []movies.WatchedMovieEntry{}
	_ = mgm.Coll(&movies.WatchedMovieEntry{}).SimpleFind(&watchedMovies, bson.M{"user_id": currentUser.ID})

	HTTPRes(c, http.StatusOK, "Watched Movies", watchedMovies)

}

func (ctl *moviesController) checkValidParameter(value string, valid []string) bool {
	for _, el := range valid {
		if el == value {
			return true
		}
	}
	return false
}
func (ctl *moviesController) ListMovies(c *gin.Context) {
	sortBy := strings.ToLower(c.Param("by"))
	direction := strings.ToLower(c.Param("direction"))

	if sortBy == "" {
		if direction != "" {
			HTTPRes(c, http.StatusBadRequest, "Validation Error", "Invalid soring method")
			return
		}
		sortBy = "name"
		direction = "desc"
	}
	if ctl.checkValidParameter(sortBy, []string{"name", "date", "rating"}) == false ||
		ctl.checkValidParameter(direction, []string{"asc", "desc"}) == false {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", "Invalid soring method")
		return
	}
	var sortByOption int8
	if direction == "asc" {
		sortByOption = 1
	} else {
		sortByOption = -1
	}

	results := []movies.MovieInfo{}
	err := mgm.Coll(&movies.Movie{}).SimpleAggregate(
		&results,
		ctl.getAggregationStages("", sortBy, sortByOption)...,
	)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}

	HTTPRes(c, http.StatusOK, "List of movies", results)

}

func (ctl *moviesController) GetMovieInfo(c *gin.Context) {
	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}

	results := []movies.MovieInfo{}
	err := mgm.Coll(&movies.Movie{}).SimpleAggregate(
		&results,
		ctl.getAggregationStages(movieId, "", 0)...,
	)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	HTTPRes(c, http.StatusOK, "List of movies", results[0])
}

func (ctl *moviesController) getAggregationStages(movieId string, sortBy string, direction int8) []interface{} {
	reviewsCollName := mgm.Coll(&movies.ReviewMovieEntry{}).Name()

	lookupStage := builder.Lookup(reviewsCollName, "_id", "movie_id", "reviews")

	countRatingsStage :=
		bson.M{operator.Set: bson.M{
			"ratingsTotal": bson.M{
				operator.Reduce: bson.M{
					"input":        "$reviews",
					"initialValue": 0,
					"in": bson.M{
						operator.Sum: bson.A{"$$value", "$$this.rating"},
					},
				},
			},
			"ratingsCount": bson.M{
				operator.Size: "$reviews",
			},
		},
		}

	averageRatingsStage :=
		bson.M{operator.Set: bson.M{
			"rating": bson.M{
				operator.Cond: bson.M{
					"if": bson.M{
						operator.Gt: bson.A{"$ratingsCount", 0},
					},
					"then": bson.M{
						operator.Divide: bson.A{"$ratingsTotal", "$ratingsCount"},
					},
					"else": 0,
				},
			},
		},
		}

	roundingStage :=
		bson.M{operator.Set: bson.M{
			"rating": bson.M{operator.Round: bson.A{"$rating", 1}},
		},
		}

	unsetStage :=
		bson.M{operator.Unset: bson.A{"ratingsCount", "ratingsTotal"}}

	var stages []interface{}
	if movieId != "" {
		movieHex, _ := primitive.ObjectIDFromHex(movieId)
		matchStage := bson.M{operator.Match: bson.M{"_id": movieHex}}
		stages = append(stages, matchStage)
	}

	stages = append(stages, lookupStage, countRatingsStage, averageRatingsStage, roundingStage, unsetStage)

	if sortBy != "" {
		sortStage := bson.M{operator.Sort: bson.M{sortBy: direction}}
		stages = append(stages, sortStage)
	}

	return stages
}

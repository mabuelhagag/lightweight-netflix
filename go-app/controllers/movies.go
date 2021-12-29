package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go-app/definitions/movies"
	"go-app/definitions/users"
	"go-app/repositories/moviesrepo"
	"go-app/repositories/usersrepo"
	"net/http"
	"time"
)

// MoviesController interface
type MoviesController interface {
	AddMovie(*gin.Context)
	UploadCover(*gin.Context)
	UpdateMovie(*gin.Context)
	DeleteMovie(c *gin.Context)
	WatchMovie(c *gin.Context)
	ReviewMovie(c *gin.Context)
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
	movie, err = ctl.mr.AddMovie(movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Failed while adding movie", err.Error())
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
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}

	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}

	currentUser := c.MustGet("user").(users.User)

	movie, err := ctl.mr.GetMovieById(movieId)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}

	watchedMovie, err := ctl.mr.DidWatchMovie(movieId, currentUser.ID.Hex())
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
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}

	err = ctl.mr.ReviewMovie(reviewEntry)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error while reviewing movie", err.Error())
		return
	}

	HTTPRes(c, http.StatusOK, "Movie reviewed", nil)

}

func (ctl *moviesController) reviewMovieInputToReviewMovieEntry(input movies.ReviewMovieInput, user users.User, movie *movies.Movie) (*movies.ReviewMovieEntry, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	return &movies.ReviewMovieEntry{
		MovieID: movie.ID,
		UserId:  user.ID,
		Rate:    input.Rate,
		Review:  input.Review,
		Time:    time.Now(),
	}, nil
}

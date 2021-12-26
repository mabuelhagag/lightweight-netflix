package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-app/definitions/movies"
	"go-app/definitions/user"
	"go-app/repositories/moviesrepo"
	"go-app/repositories/userrepo"
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
	ur userrepo.Repo
}

// NewMoviesController instantiates User Controller
func NewMoviesController(br moviesrepo.Repo, us userrepo.Repo) MoviesController {
	return &moviesController{mr: br, ur: us}
}

func (ctl *moviesController) AddMovie(c *gin.Context) {
	var movieInput movies.AddMovieInput
	if err := c.ShouldBindJSON(&movieInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}
	email := c.MustGet("email").(string)
	movie, err := ctl.inputToMovie(movieInput, email)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
	}
	movie, err = ctl.mr.AddMovie(movie)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Failed while adding movie", err.Error())
	}
	HTTPRes(c, http.StatusOK, "Movie added", movie)
}

func (ctl *moviesController) inputToMovie(input movies.AddMovieInput, email string) (*movies.Movie, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	currentUser, err := ctl.ur.GetUser(email)
	if err != nil {
		return nil, err
	}

	return &movies.Movie{
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
		AddedBy:     currentUser.ID,
	}, nil
}

func (ctl *moviesController) UploadCover(c *gin.Context) {
	var uploadCoverInput movies.UploadCoverInput
	if err := c.ShouldBind(&uploadCoverInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}
	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}
	_, err := ctl.mr.GetMovieById(movieId)
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
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}

	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}

	email := c.MustGet("email").(string)
	currentUser, err := ctl.ur.GetUser(email)
	movie, err := ctl.mr.GetMovieById(movieId)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	if movie.AddedBy != currentUser.ID {
		HTTPRes(c, http.StatusForbidden, "Error updating movie info", "Movie is not owned by current user")
		return
	}
	movieUpdate, err := ctl.updateMovieInputToMovie(movieInput)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}
	movie, err = ctl.mr.UpdateMovie(movieUpdate, movieId)

	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Failed while adding movie", err.Error())
		return
	}
	HTTPRes(c, http.StatusOK, "Movie Updated", movie)
}

func (ctl *moviesController) updateMovieInputToMovie(input movies.UpdateMovieInput) (*movies.Movie, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	return &movies.Movie{
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
	}, nil
}

func (ctl *moviesController) DeleteMovie(c *gin.Context) {
	movieId := c.Param("id")
	if movieId == "" {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", "Movie ID not provided")
		return
	}

	email := c.MustGet("email").(string)
	currentUser, err := ctl.ur.GetUser(email)
	movie, err := ctl.mr.GetMovieById(movieId)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting movie info", err.Error())
		return
	}
	if movie.AddedBy != currentUser.ID {
		HTTPRes(c, http.StatusForbidden, "Error updating movie info", "Movie is not owned by current user")
		return
	}
	err = ctl.mr.DeleteMovie(movieId)
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

	email := c.MustGet("email").(string)
	currentUser, err := ctl.ur.GetUser(email)
	movie, err := ctl.mr.GetMovieById(movieId)
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

	email := c.MustGet("email").(string)
	currentUser, err := ctl.ur.GetUser(email)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Error getting user", err.Error())
		return
	}

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

func (ctl *moviesController) reviewMovieInputToReviewMovieEntry(input movies.ReviewMovieInput, user *user.User, movie *movies.Movie) (*movies.ReviewMovieEntry, error) {
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

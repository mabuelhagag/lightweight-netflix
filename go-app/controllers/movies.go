package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	moviesdefinition "go-app/definitions/movies"
	"go-app/repositories/moviesrepo"
	"go-app/repositories/userrepo"
	"net/http"
)

// MoviesController interface
type MoviesController interface {
	AddMovie(*gin.Context)
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
	var movieInput moviesdefinition.AddMovieInput
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

func (ctl *moviesController) inputToMovie(input moviesdefinition.AddMovieInput, email string) (*moviesdefinition.Movie, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	user, err := ctl.ur.GetUser(email)
	if err != nil {
		return nil, err
	}

	return &moviesdefinition.Movie{
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
		AddedBy:     user.ID,
	}, nil
}

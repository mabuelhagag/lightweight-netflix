package controllers

import (
	"context"
	"go-app/repositories/userrepo"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/mold/v4/modifiers"
)

// UserOutput represents HTTP Response Body structure
type UserOutput struct {
	Name  string `json:"name"`
	Pages uint   `json:"pages"`
}

// UserInput represents createUser body format
type UserInput struct {
	FullName             string `json:"full_name" mod:"trim" binding:"required"`
	Age                  uint8  `json:"age" binding:"required,gt=0,lt=100"`
	Email                string `json:"email" mod:"trim,lcase" binding:"required"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

var conform = modifiers.New()

// UserController interface
type UserController interface {
	RegisterUser(*gin.Context)
}

type userController struct {
	br userrepo.Repo
}

// NewUserController instantiates User Controller
func NewUserController(br userrepo.Repo) UserController {
	return &userController{br: br}
}

func (ctl *userController) RegisterUser(c *gin.Context) {
	// Read user input
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}
	b, err := ctl.inputToUser(userInput, c)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", nil)
		return
	}

	// Create user
	// If an Error Occurs while creating return the error
	if _, err := ctl.br.CreateUser(b); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// If user is successfully created return a structured Response
	//userOutput := ctl.mapToUserOutput(&b)
	HTTPRes(c, http.StatusOK, "User Registered", nil)
}

// Private Methods
func (ctl *userController) inputToUser(input UserInput, c *gin.Context) (*userrepo.User, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return nil, err
	}

	return &userrepo.User{
		FullName: input.FullName,
		Age:      input.Age,
		Email:    input.Email,
		Password: input.Password,
	}, nil
}

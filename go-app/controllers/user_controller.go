package controllers

import (
	"go-app/domain/user"
	"go-app/services/userservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserOutput represents HTTP Response Body structure
type UserOutput struct {
	Name  string `json:"name"`
	Pages uint   `json:"pages"`
}

// UserInput represents createUser body format
type UserInput struct {
	Name  string `json:"name"`
	Pages uint   `json:"pages"`
}

// UserController interface
type UserController interface {
	RegisterUser(*gin.Context)
}

type userController struct {
	bs userservice.UserService
}

// NewUserController instantiates User Controller
func NewUserController(bs userservice.UserService) UserController {
	return &userController{bs: bs}
}

func (ctl *userController) RegisterUser(c *gin.Context) {
	// Read user input
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	b := ctl.inputToUser(userInput)

	// Create user
	// If an Error Occurs while creating return the error
	if _, err := ctl.bs.CreateUser(&b); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// If user is successfully created return a structured Response
	userOutput := ctl.mapToUserOutput(&b)
	HTTPRes(c, http.StatusOK, "User Registered", userOutput)
}

// Private Methods
func (ctl *userController) inputToUser(input UserInput) user.User {
	return user.User{
		Name:  input.Name,
		Pages: input.Pages,
	}
}
func (ctl *userController) mapToUserOutput(b *user.User) *UserOutput {
	return &UserOutput{
		Name:  b.Name,
		Pages: b.Pages,
	}
}

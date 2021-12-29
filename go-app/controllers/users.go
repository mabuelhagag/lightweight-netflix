package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	userdefinition "go-app/definitions/user"
	"go-app/repositories/userrepo"
	"log"
	"net/http"
)

// UserController interface
type UserController interface {
	RegisterUser(*gin.Context)
	LoginUser(*gin.Context)
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
	var userInput userdefinition.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	b, err := ctl.inputToUser(userInput)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	// Create user
	// If an Error Occurs while creating return the error
	if _, err := ctl.br.CreateUser(b); err != nil {
		HTTPRes(c, http.StatusInternalServerError, "Failed while registering user", err.Error())
		return
	}

	HTTPRes(c, http.StatusOK, "User Registered", nil)
}

// Private Methods
func (ctl *userController) inputToUser(input userdefinition.UserInput) (*userdefinition.User, error) {
	if err := conform.Struct(context.Background(), &input); err != nil {
		return nil, err
	}

	return &userdefinition.User{
		FullName: input.FullName,
		Age:      input.Age,
		Email:    input.Email,
		Password: input.Password,
	}, nil
}

func (ctl *userController) LoginUser(c *gin.Context) {
	var loginInfoInput userdefinition.LoginInfoInput
	if err := c.ShouldBindJSON(&loginInfoInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}

	if err := conform.Struct(context.Background(), &loginInfoInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, "Error Validation", err.Error())
		return
	}

	if err := ctl.br.CheckPassword(&loginInfoInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	signedToken, err := userdefinition.AppJwtWrapper.GenerateToken(loginInfoInput.Email)
	if err != nil {
		log.Println(err)
		HTTPRes(c, http.StatusInternalServerError, "error signing token", nil)
		return
	}
	tokenResponse := userdefinition.LoginInfoOutput{
		Token: signedToken,
	}

	HTTPRes(c, http.StatusBadRequest, "User Authorized", tokenResponse)
	return
}

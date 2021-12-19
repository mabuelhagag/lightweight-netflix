package api

import (
	"errors"
	"strings"
)

// UserService contains the methods of the user service
type UserService interface {
	New(user NewUserRequest) error
}

// User repository is what lets our service do db operations without knowing anything about the implementation
type UserRepository interface {
	CreateUser(NewUserRequest) error
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{storage: userRepo}
}

func (u *userService) New(user NewUserRequest) error {
	// sanitize fields
	//fullName := strings.TrimSpace(user.FullName)
	//email := strings.TrimSpace(strings.ToLower(user.Email))
	// age,_ := strconv.ParseUint(user.Age)
	// do some basic validations

	if user.Email == "" {
		return errors.New("user service - email required")
	}

	if user.FullName == "" {
		return errors.New("user service - name required")
	}

	// do some basic normalisation
	user.FullName = strings.ToLower(user.FullName)
	user.Email = strings.TrimSpace(user.Email)

	err := u.storage.CreateUser(user)

	if err != nil {
		return err
	}

	return nil
}
